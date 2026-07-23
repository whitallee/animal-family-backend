package species

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/whitallee/animal-family-backend/types"
)

type Store struct {
	db             *sql.DB
	openAIKey      string
	s3AssetsBucket string
	awsRegion      string
}

func NewStore(db *sql.DB, openAIKey, s3AssetsBucket, awsRegion string) *Store {
	return &Store{
		db:             db,
		openAIKey:      openAIKey,
		s3AssetsBucket: s3AssetsBucket,
		awsRegion:      awsRegion,
	}
}

func (s *Store) CreateSpecies(species types.Species) error {
	_, err := s.db.Exec(`INSERT INTO "species" ("comName", "sciName", "speciesDesc", "image", "habitatId", "baskTemp", "diet", "sociality", "lifespan", "size", "weight", "conservationStatus", "extraCare") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`, species.ComName, species.SciName, species.SpeciesDesc, species.Image, species.HabitatId, species.BaskTemp, species.Diet, species.Sociality, species.Lifespan, species.Size, species.Weight, species.ConservationStatus, species.ExtraCare)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateSpecies(species types.Species) error {
	_, err := s.db.Exec(`UPDATE "species"
						SET "comName" = $1, "sciName" = $2, "speciesDesc" = $3, "image" = $4, "habitatId" = $5, "baskTemp" = $6, "diet" = $7, "sociality" = $8, "lifespan" = $9, "size" = $10, "weight" = $11, "conservationStatus" = $12, "extraCare" = $13
						WHERE "speciesId" = $14`, species.ComName, species.SciName, species.SpeciesDesc, species.Image, species.HabitatId, species.BaskTemp, species.Diet, species.Sociality, species.Lifespan, species.Size, species.Weight, species.ConservationStatus, species.ExtraCare, species.SpeciesID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetSpecies() ([]*types.Species, error) {
	rows, err := s.db.Query(`SELECT * FROM "species"`)
	if err != nil {
		return nil, err
	}

	species := make([]*types.Species, 0)
	for rows.Next() {
		s, err := scanRowsIntoSpecies(rows)
		if err != nil {
			return nil, err
		}

		species = append(species, s)
	}
	return species, nil
}

func (s *Store) GetSpeciesByComName(comName string) (*types.Species, error) {
	rows, err := s.db.Query(`SELECT * FROM "species" WHERE "comName" = $1`, comName)
	if err != nil {
		return nil, err
	}

	species := new(types.Species)
	for rows.Next() {
		species, err = scanRowsIntoSpecies(rows)
		if err != nil {
			return nil, err
		}
	}

	if species.SpeciesID == 0 {
		return nil, fmt.Errorf("species not found")
	}

	return species, nil
}

func (s *Store) GetSpeciesBySciName(sciName string) (*types.Species, error) {
	rows, err := s.db.Query(`SELECT * FROM "species" WHERE "sciName" = $1`, sciName)
	if err != nil {
		return nil, err
	}

	species := new(types.Species)
	for rows.Next() {
		species, err = scanRowsIntoSpecies(rows)
		if err != nil {
			return nil, err
		}
	}

	if species.SpeciesID == 0 {
		return nil, fmt.Errorf("species not found")
	}

	return species, nil
}

func (s *Store) GetSpeciesById(speciesId int) (*types.Species, error) {
	rows, err := s.db.Query(`SELECT * FROM "species" WHERE "speciesId" = $1`, speciesId)
	if err != nil {
		return nil, err
	}

	species := new(types.Species)
	for rows.Next() {
		species, err = scanRowsIntoSpecies(rows)
		if err != nil {
			return nil, err
		}
	}

	if species.SpeciesID == 0 {
		return nil, fmt.Errorf("species not found")
	}

	return species, nil
}

func (s *Store) DeleteSpeciesById(speciesId int) error {
	// begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// update animals species to "No Species"
	_, err = tx.Exec(`UPDATE "animals" SET "speciesId" = 0 WHERE "speciesid" = $1`, speciesId)
	if err != nil {
		return err
	}

	// delete species
	_, err = tx.Exec(`DELETE FROM "species" WHERE "speciesId" = $1`, speciesId)
	if err != nil {
		return err
	}

	// commit changes
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetSpeciesByNameCaseInsensitive(name string) (*types.Species, error) {
	rows, err := s.db.Query(`SELECT * FROM "species" WHERE LOWER("comName") = LOWER($1) OR LOWER("sciName") = LOWER($1) LIMIT 1`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	species := new(types.Species)
	for rows.Next() {
		species, err = scanRowsIntoSpecies(rows)
		if err != nil {
			return nil, err
		}
	}

	if species.SpeciesID == 0 {
		return nil, fmt.Errorf("species not found")
	}

	return species, nil
}

func (s *Store) GenerateSpeciesFromName(name string) (*types.Species, error) {
	// Fetch habitats to give the LLM context for picking habitatId
	rows, err := s.db.Query(`SELECT "habitatId", "habitatName" FROM "habitats"`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch habitats: %w", err)
	}
	var habitatLines []string
	for rows.Next() {
		var id int
		var habitatName string
		if scanErr := rows.Scan(&id, &habitatName); scanErr != nil {
			rows.Close()
			return nil, scanErr
		}
		habitatLines = append(habitatLines, fmt.Sprintf("ID %d: %s", id, habitatName))
	}
	rows.Close()
	habitatList := strings.Join(habitatLines, "\n")

	// Two-call LLM flow: natural language → structured JSON
	llmData, err := generateSpeciesData(s.openAIKey, name, habitatList)
	if err != nil {
		return nil, err
	}

	imageFilename := speciesImageFilename(llmData.ComName)
	placeholderURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/species/no-species.png", s.s3AssetsBucket, s.awsRegion)

	// Insert with placeholder image and get back the new ID
	var speciesId int
	err = s.db.QueryRow(`
		INSERT INTO "species" ("comName", "sciName", "speciesDesc", "image", "habitatId", "baskTemp", "diet", "sociality", "lifespan", "size", "weight", "conservationStatus", "extraCare")
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING "speciesId"`,
		llmData.ComName, llmData.SciName, llmData.SpeciesDesc, placeholderURL, llmData.HabitatId,
		llmData.BaskTemp, llmData.Diet, llmData.Sociality, llmData.Lifespan, llmData.Size,
		llmData.Weight, llmData.ConservationStatus, llmData.ExtraCare,
	).Scan(&speciesId)
	if err != nil {
		return nil, fmt.Errorf("failed to insert species: %w", err)
	}

	theme := ImageThemeTerrestrial
	if ImageTheme(llmData.ImageTheme) == ImageThemeAquatic {
		theme = ImageThemeAquatic
	}

	// Non-blocking: generate image, upload to S3, update DB record
	go generateAndUploadSpeciesImage(s, speciesId, llmData.ComName, imageFilename, theme)

	return &types.Species{
		SpeciesID:          speciesId,
		ComName:            llmData.ComName,
		SciName:            llmData.SciName,
		SpeciesDesc:        llmData.SpeciesDesc,
		Image:              placeholderURL,
		HabitatId:          llmData.HabitatId,
		BaskTemp:           llmData.BaskTemp,
		Diet:               llmData.Diet,
		Sociality:          llmData.Sociality,
		Lifespan:           llmData.Lifespan,
		Size:               llmData.Size,
		Weight:             llmData.Weight,
		ConservationStatus: llmData.ConservationStatus,
		ExtraCare:          llmData.ExtraCare,
	}, nil
}

func (s *Store) UpdateSpeciesImage(speciesId int, imageURL string) error {
	_, err := s.db.Exec(`UPDATE "species" SET "image" = $1 WHERE "speciesId" = $2`, imageURL, speciesId)
	return err
}

func scanRowsIntoSpecies(rows *sql.Rows) (*types.Species, error) {
	species := new(types.Species)

	err := rows.Scan(
		&species.SpeciesID,
		&species.ComName,
		&species.SciName,
		&species.Image,
		&species.SpeciesDesc,
		&species.HabitatId,
		&species.BaskTemp,
		&species.Diet,
		&species.Sociality,
		&species.Lifespan,
		&species.Size,
		&species.Weight,
		&species.ConservationStatus,
		&species.ExtraCare,
	)
	if err != nil {
		return nil, err
	}

	return species, nil
}
