package species

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ImageTheme controls the background style used in the DALL-E prompt
type ImageTheme string

const (
	ImageThemeAquatic     ImageTheme = "aquatic"
	ImageThemeTerrestrial ImageTheme = "terrestrial"
)

func (t ImageTheme) backgroundPrompt() string {
	if t == ImageThemeAquatic {
		return "Soft aquamarine paper-textured background (#69C7CB style), subtle grain texture, slight color variation, no plants, bubbles, rocks, or scenery."
	}
	return "Warm tan paper-textured background (#E4BE72 style), subtle grain texture, slight color variation, no objects or scenery."
}

const dalleBasePrompt = `Create a square (1:1) fallback illustration of a {ANIMAL_NAME}.

Style:
- Flat, stylized digital illustration
- Mid-century modern inspired
- Simple geometric shapes
- Clean outlines
- Minimal details
- Soft paper-like grain texture throughout the image
- Friendly and approachable appearance
- Designed as an app fallback image/icon, not realistic wildlife art
- No photorealism
- No dramatic lighting
- No shadows except a subtle grounding shadow if needed
- No background scenery
- No plants, rocks, branches, or habitat elements
- No text

Composition:
- Animal centered in frame
- Occupies roughly 60–70% of the canvas
- Entire animal visible
- Facing left at a slight 3/4 angle
- Consistent composition with a species icon library
- Generous negative space around the animal

Color:
- Use the natural colors and recognizable markings of a {ANIMAL_NAME}
- Simplify markings into bold graphic shapes
- Preserve distinctive identifying features

Background:
{BACKGROUND_DESCRIPTION}

Rendering:
- Flat color blocks
- Subtle texture overlay
- Rounded friendly shapes
- Consistent with a modern animal encyclopedia or premium pet-care app
- Vector-like appearance
- Avoid realistic scales, fur strands, reflections, or photographic detail

The final image should feel like part of a cohesive icon set alongside stylized illustrations of leopard geckos, tree frogs, gouramis, turtles, chameleons, and betta fish.`

const speciesInfoSystemPrompt = `You are an expert zoologist and animal care specialist. Provide detailed, accurate information about animal species including their biology, natural habitat, temperature requirements, diet, social behavior, lifespan, size, conservation status, and captive care requirements.`

const speciesInfoUserPromptFmt = `Provide comprehensive information about the following species: %s

Include details covering:
- General description and distinguishing characteristics
- Natural habitat and environment
- Temperature and environmental requirements
- Diet in both the wild and captivity
- Social behavior (solitary, schooling, etc.) and ideal group sizes if applicable
- Typical lifespan in captivity
- Adult size and weight
- IUCN conservation status
- Special care requirements for keeping this species in captivity

Available habitat categories (note which best fits this species):
%s`

const speciesJSONSystemPrompt = `You are a data formatting assistant. Convert species information into a structured JSON object. Return ONLY valid JSON with no additional text or markdown.`

const speciesJSONUserPromptFmt = `Convert the following species information into a JSON object with these exact fields:
- comName: common name
- sciName: scientific name in "Genus species" format
- speciesDesc: 2-4 sentence description suitable for a pet care app (max 1500 chars)
- habitatId: integer ID chosen from the habitat list below
- baskTemp: temperature range string (e.g. "72-78°F (22-26°C)")
- diet: diet description (max 1500 chars)
- sociality: social behavior description including group size recommendations (max 1500 chars)
- lifespan: lifespan string (e.g. "3-5 years")
- size: adult size string (e.g. "1 inch (2.5 cm)")
- weight: adult weight string (e.g. "Less than 1 gram")
- conservationStatus: IUCN status (e.g. "Least Concern", "Near Threatened")
- extraCare: special care requirements (max 1500 chars)
- imageTheme: "aquatic" for fish and fully aquatic animals, "terrestrial" for land and semi-aquatic animals

Available habitat IDs:
%s

Species information to convert:
%s`

// Internal types for OpenAI API communication

type openAIMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRespFmt struct {
	Type string `json:"type"`
}

type openAIChatReq struct {
	Model          string         `json:"model"`
	Messages       []openAIMsg    `json:"messages"`
	ResponseFormat *openAIRespFmt `json:"response_format,omitempty"`
}

type openAIChatResp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type openAIImageReq struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	N       int    `json:"n"`
	Size    string `json:"size"`
	Quality string `json:"quality"`
}

type openAIImageResp struct {
	Data []struct {
		B64JSON string `json:"b64_json"`
	} `json:"data"`
}

type llmSpeciesJSON struct {
	ComName            string `json:"comName"`
	SciName            string `json:"sciName"`
	SpeciesDesc        string `json:"speciesDesc"`
	HabitatId          int    `json:"habitatId"`
	BaskTemp           string `json:"baskTemp"`
	Diet               string `json:"diet"`
	Sociality          string `json:"sociality"`
	Lifespan           string `json:"lifespan"`
	Size               string `json:"size"`
	Weight             string `json:"weight"`
	ConservationStatus string `json:"conservationStatus"`
	ExtraCare          string `json:"extraCare"`
	ImageTheme         string `json:"imageTheme"`
}

// generateSpeciesData runs the two-call LLM flow and returns parsed species data
func generateSpeciesData(apiKey, name, habitatList string) (*llmSpeciesJSON, error) {
	naturalInfo, err := callOpenAIChat(
		apiKey,
		speciesInfoSystemPrompt,
		fmt.Sprintf(speciesInfoUserPromptFmt, name, habitatList),
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate species info: %w", err)
	}

	jsonStr, err := callOpenAIChat(
		apiKey,
		speciesJSONSystemPrompt,
		fmt.Sprintf(speciesJSONUserPromptFmt, habitatList, naturalInfo),
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to structure species data: %w", err)
	}

	var data llmSpeciesJSON
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("failed to parse LLM JSON response: %w", err)
	}

	return &data, nil
}

// callOpenAIChat sends a chat completion request to OpenAI
func callOpenAIChat(apiKey, systemPrompt, userPrompt string, jsonMode bool) (string, error) {
	reqBody := openAIChatReq{
		Model: "gpt-4o",
		Messages: []openAIMsg{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}
	if jsonMode {
		reqBody.ResponseFormat = &openAIRespFmt{Type: "json_object"}
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai chat error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result openAIChatResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}

	return result.Choices[0].Message.Content, nil
}

// generateAndUploadSpeciesImage runs in a goroutine: generates image via DALL-E,
// uploads to S3, and updates the species DB record with the real URL
func generateAndUploadSpeciesImage(store *Store, speciesId int, comName, imageFilename string, theme ImageTheme) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("species image generation panic for %s: %v", comName, r)
		}
	}()

	prompt := strings.ReplaceAll(dalleBasePrompt, "{ANIMAL_NAME}", comName)
	prompt = strings.ReplaceAll(prompt, "{BACKGROUND_DESCRIPTION}", theme.backgroundPrompt())

	imageBytes, err := callOpenAIImageGen(store.openAIKey, prompt)
	if err != nil {
		log.Printf("failed to generate image for %s: %v", comName, err)
		return
	}

	s3Key := "species/" + imageFilename
	if err := uploadToS3(context.Background(), store.awsRegion, store.s3AssetsBucket, s3Key, imageBytes, "image/png"); err != nil {
		log.Printf("failed to upload image to S3 for %s: %v", comName, err)
		return
	}

	s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", store.s3AssetsBucket, store.awsRegion, s3Key)
	if err := store.UpdateSpeciesImage(speciesId, s3URL); err != nil {
		log.Printf("failed to update species image URL for %s: %v", comName, err)
		return
	}

	log.Printf("image generated and uploaded for %s: %s", comName, s3URL)
}

// callOpenAIImageGen calls gpt-image-2 and returns the decoded image bytes
func callOpenAIImageGen(apiKey, prompt string) ([]byte, error) {
	reqBody := openAIImageReq{
		Model:   "gpt-image-2",
		Prompt:  prompt,
		N:       1,
		Size:    "1024x1024",
		Quality: "medium",
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/images/generations", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai image error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result openAIImageResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("openai returned no image data")
	}

	imageBytes, err := base64.StdEncoding.DecodeString(result.Data[0].B64JSON)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 image: %w", err)
	}

	return imageBytes, nil
}

// uploadToS3 uploads bytes to the given S3 bucket and key
func uploadToS3(ctx context.Context, region, bucket, key string, data []byte, contentType string) error {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(len(data))),
	})
	return err
}

// speciesImageFilename converts a common name to a kebab-case filename
func speciesImageFilename(comName string) string {
	name := strings.ToLower(strings.TrimSpace(comName))
	name = strings.ReplaceAll(name, " ", "-")
	var result strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String() + ".png"
}
