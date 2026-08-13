package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type parseResult struct {
	id         int
	err        error
	handlerRan bool
	status     int
}

// parseIDFrom routes a request through mux so that mux.Vars is populated the
// same way it is in production.
func parseIDFrom(t *testing.T, url string) parseResult {
	t.Helper()

	var res parseResult

	router := mux.NewRouter()
	router.HandleFunc("/habitats/{id}", func(w http.ResponseWriter, r *http.Request) {
		res.handlerRan = true
		res.id, res.err = ParseIDParam(r, "id")
	})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, url, nil))
	res.status = rr.Code

	return res
}

func TestParseIDParamValid(t *testing.T) {
	res := parseIDFrom(t, "/habitats/42")
	if res.err != nil {
		t.Fatalf("unexpected error: %v", res.err)
	}
	if res.id != 42 {
		t.Errorf("expected 42, got %d", res.id)
	}
}

func TestParseIDParamRejectsInvalid(t *testing.T) {
	cases := map[string]string{
		"non-numeric": "/habitats/abc",
		"zero":        "/habitats/0",
		"negative":    "/habitats/-5",
		"decimal":     "/habitats/1.5",
	}

	for name, url := range cases {
		t.Run(name, func(t *testing.T) {
			res := parseIDFrom(t, url)
			if !res.handlerRan {
				t.Fatalf("expected %s to reach the handler", url)
			}
			if res.err == nil {
				t.Errorf("expected an error for %s, got none (parsed id %d)", url, res.id)
			}
		})
	}
}

// A missing ID segment is rejected by mux before the handler is reached, so
// ParseIDParam is never consulted. Asserting the 404 keeps that guarantee
// visible: it is why handlers can treat a reached-handler request as having
// some value present for {id}.
func TestParseIDParamMissingSegmentIsNotRouted(t *testing.T) {
	res := parseIDFrom(t, "/habitats/")
	if res.handlerRan {
		t.Error("expected /habitats/ not to match the /habitats/{id} route")
	}
	if res.status != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, res.status)
	}
}

func TestParseIDParamMissingParam(t *testing.T) {
	// A request that never went through mux has no vars at all.
	r := httptest.NewRequest(http.MethodGet, "/habitats", nil)
	if _, err := ParseIDParam(r, "id"); err == nil {
		t.Error("expected an error when the path parameter is absent")
	}
}
