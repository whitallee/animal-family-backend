// Command specgen turns the Swagger 2.0 document produced by swag into the
// OpenAPI 3.0 document the frontend's client generator consumes.
//
// It does two things:
//
//  1. Strips the Go package prefix from schema names, so "types.Habitat"
//     becomes "Habitat". swag names definitions after the fully qualified Go
//     type and its @name override only works via a trailing line comment on the
//     type (it reads TypeSpec.Comment, never GenDecl.Doc), which would mean
//     awkward comment placement on every exported type. Renaming here is one
//     deterministic step that covers every type instead.
//
//  2. Converts Swagger 2.0 to OpenAPI 3.0. swag emits 2.0 only; orval needs 3.x.
//
// Run via `make spec`, which regenerates docs/swagger.json first.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
)

const definitionsPrefix = "#/definitions/"

func main() {
	in := flag.String("in", "docs/swagger.json", "Swagger 2.0 input produced by swag")
	out := flag.String("out", "docs/openapi.json", "OpenAPI 3.0 output")
	flag.Parse()

	raw, err := os.ReadFile(*in)
	if err != nil {
		log.Fatalf("read %s: %v (run `make spec` rather than this tool directly)", *in, err)
	}

	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		log.Fatalf("parse %s: %v", *in, err)
	}

	renames, err := buildRenames(doc)
	if err != nil {
		log.Fatalf("schema names: %v", err)
	}
	applyRenames(doc, renames)

	log.Printf("marked %d response schemas as fully required", markResponseFieldsRequired(doc))

	normalized, err := json.Marshal(doc)
	if err != nil {
		log.Fatalf("re-encode: %v", err)
	}

	var v2 openapi2.T
	if err := json.Unmarshal(normalized, &v2); err != nil {
		log.Fatalf("decode swagger 2.0: %v", err)
	}

	v3, err := openapi2conv.ToV3(&v2)
	if err != nil {
		log.Fatalf("convert to openapi 3.0: %v", err)
	}

	// The 2.0 -> 3.0 conversion only produces a `servers` entry when `host` is
	// set, which it deliberately is not (the API is served from different hosts
	// per environment). Without this the spec loses the /api/v2 prefix entirely
	// and generated clients would request /habitats instead of
	// /api/v2/habitats.
	if len(v3.Servers) == 0 && v2.BasePath != "" {
		v3.Servers = openapi3.Servers{{
			URL:         v2.BasePath,
			Description: "Relative to the API host, configured per environment.",
		}}
	}

	encoded, err := json.MarshalIndent(v3, "", "  ")
	if err != nil {
		log.Fatalf("encode openapi 3.0: %v", err)
	}
	encoded = append(encoded, '\n')

	if err := os.WriteFile(*out, encoded, 0o644); err != nil {
		log.Fatalf("write %s: %v", *out, err)
	}

	log.Printf("wrote %s (%d paths, %d schemas)", *out, len(v3.Paths.Map()), len(v3.Components.Schemas))
}

// buildRenames maps each definition name to its package-qualifier-stripped
// form. It fails rather than silently dropping a schema if two packages would
// collapse onto the same name.
func buildRenames(doc map[string]any) (map[string]string, error) {
	defs, ok := doc["definitions"].(map[string]any)
	if !ok {
		return map[string]string{}, nil
	}

	renames := make(map[string]string, len(defs))
	claimed := make(map[string]string, len(defs))

	names := make([]string, 0, len(defs))
	for name := range defs {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		short := name
		if idx := strings.LastIndex(name, "."); idx != -1 {
			short = name[idx+1:]
		}
		if short == "" {
			continue
		}
		if prior, taken := claimed[short]; taken {
			return nil, fmt.Errorf("%q and %q both shorten to %q; give one an explicit distinct Go type name", prior, name, short)
		}
		claimed[short] = name
		renames[name] = short
	}

	return renames, nil
}

// applyRenames rewrites definition keys and every $ref that points at them.
// It walks the decoded document rather than doing text substitution so that
// descriptions mentioning a type name are left alone.
func applyRenames(doc map[string]any, renames map[string]string) {
	if defs, ok := doc["definitions"].(map[string]any); ok {
		renamed := make(map[string]any, len(defs))
		for name, schema := range defs {
			if short, ok := renames[name]; ok {
				renamed[short] = schema
			} else {
				renamed[name] = schema
			}
		}
		doc["definitions"] = renamed
	}

	rewriteRefs(doc, renames)
}

// markResponseFieldsRequired marks every property of a response-only schema as
// required, and returns how many schemas it changed.
//
// swag derives `required` from `validate:"required"` tags, which only request
// payloads carry. Response models have no such tags, so without this every
// field of every response type generates as optional in TypeScript — callers
// would have to null-check fields the server always sends.
//
// This is accurate rather than optimistic: no json tag in types/ uses
// omitempty, so Go serializes every exported field on every response.
//
// Schemas used by both a request and a response are skipped, since their
// existing `required` set is meaningful for the request side.
//
// NOTE: when nullable fields reach the spec (Animal's *string and
// sql.NullString fields, arriving with enclosures/animals), required will still
// be correct — the key is always present — but those properties will also need
// marking nullable so the generated types admit null.
func markResponseFieldsRequired(doc map[string]any) int {
	defs, ok := doc["definitions"].(map[string]any)
	if !ok {
		return 0
	}

	inResponses := map[string]bool{}
	inRequests := map[string]bool{}

	paths, _ := doc["paths"].(map[string]any)
	for _, pathItem := range paths {
		operations, ok := pathItem.(map[string]any)
		if !ok {
			continue
		}
		for _, operation := range operations {
			op, ok := operation.(map[string]any)
			if !ok {
				continue
			}

			if responses, ok := op["responses"].(map[string]any); ok {
				for _, response := range responses {
					if r, ok := response.(map[string]any); ok {
						collectRefs(r["schema"], defs, inResponses)
					}
				}
			}

			if params, ok := op["parameters"].([]any); ok {
				for _, param := range params {
					p, ok := param.(map[string]any)
					if !ok || p["in"] != "body" {
						continue
					}
					collectRefs(p["schema"], defs, inRequests)
				}
			}
		}
	}

	changed := 0
	for name := range inResponses {
		if inRequests[name] {
			continue
		}

		schema, ok := defs[name].(map[string]any)
		if !ok {
			continue
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok || len(properties) == 0 {
			continue
		}
		if _, already := schema["required"]; already {
			continue
		}

		required := make([]string, 0, len(properties))
		for property := range properties {
			required = append(required, property)
		}
		sort.Strings(required)

		schema["required"] = required
		changed++
	}

	return changed
}

// collectRefs records every definition reachable from node, following $refs
// through nested schemas so that a type referenced only via another response
// type is still covered.
func collectRefs(node any, defs map[string]any, found map[string]bool) {
	switch typed := node.(type) {
	case map[string]any:
		for key, value := range typed {
			if key == "$ref" {
				ref, ok := value.(string)
				if !ok || !strings.HasPrefix(ref, definitionsPrefix) {
					continue
				}
				name := strings.TrimPrefix(ref, definitionsPrefix)
				if found[name] {
					continue
				}
				found[name] = true
				collectRefs(defs[name], defs, found)
				continue
			}
			collectRefs(value, defs, found)
		}
	case []any:
		for _, item := range typed {
			collectRefs(item, defs, found)
		}
	}
}

func rewriteRefs(node any, renames map[string]string) {
	switch typed := node.(type) {
	case map[string]any:
		for key, value := range typed {
			if key == "$ref" {
				if ref, ok := value.(string); ok {
					if short, found := renames[strings.TrimPrefix(ref, definitionsPrefix)]; found && strings.HasPrefix(ref, definitionsPrefix) {
						typed[key] = definitionsPrefix + short
					}
				}
				continue
			}
			rewriteRefs(value, renames)
		}
	case []any:
		for _, item := range typed {
			rewriteRefs(item, renames)
		}
	}
}
