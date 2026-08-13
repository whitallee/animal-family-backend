package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ParseIDParam reads a positive integer ID from a mux path parameter.
//
// v2 routes take resource IDs as path parameters rather than in a JSON body,
// so this replaces the `ParseJSON` + `Validate.Struct` pair that v1 used on
// `<X>IdPayload` structs. Keeping the positive-integer check here preserves the
// `validate:"required,min=0"` coverage those payloads provided.
func ParseIDParam(r *http.Request, name string) (int, error) {
	raw, ok := mux.Vars(r)[name]
	if !ok {
		return 0, fmt.Errorf("missing %s path parameter", name)
	}

	id, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: must be an integer", name)
	}

	if id < 1 {
		return 0, fmt.Errorf("invalid %s: must be a positive integer", name)
	}

	return id, nil
}
