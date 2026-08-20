package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// serve runs handler against a real http.Server rather than httptest.Recorder,
// because the Recorder does not enforce HTTP body rules — the 204 behaviour
// these tests cover is invisible to it.
func serve(t *testing.T, handler http.HandlerFunc) (*http.Response, string) {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	response, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	t.Cleanup(func() {
		if err := response.Body.Close(); err != nil {
			t.Errorf("close body: %v", err)
		}
	})

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	return response, string(body)
}

func TestWriteJSONEncodesValue(t *testing.T) {
	response, body := serve(t, func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", response.StatusCode)
	}
	if body != "{\"status\":\"ok\"}\n" {
		t.Errorf("unexpected body %q", body)
	}
	if got := response.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("expected application/json, got %q", got)
	}
}

// 204 must not carry a body. Encoding one made net/http return an error on
// every successful update and delete, which is what the discarded `_ =` at each
// call site was hiding.
func TestWriteJSONNoContentSendsNoBody(t *testing.T) {
	response, body := serve(t, func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusNoContent, nil)
	})

	if response.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", response.StatusCode)
	}
	if body != "" {
		t.Errorf("expected an empty body for 204, got %q", body)
	}
}

// The v1 frontend calls res.json() on creates without checking the status, so
// an empty 201 body would throw "Unexpected end of JSON input" and report a
// failure for a create that actually succeeded. Keep writing `null` until those
// callers are migrated to v2.
func TestWriteJSONCreatedKeepsNullBodyForV1Clients(t *testing.T) {
	response, body := serve(t, func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusCreated, nil)
	})

	if response.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", response.StatusCode)
	}
	if body != "null\n" {
		t.Errorf("expected \"null\\n\" for 201, got %q", body)
	}
}

// v2 creates document a 201 with no content, so the generated client types the
// body as void. WriteStatus is what makes that true on the wire.
func TestWriteStatusSendsNoBody(t *testing.T) {
	response, body := serve(t, func(w http.ResponseWriter, r *http.Request) {
		WriteStatus(w, http.StatusCreated)
	})

	if response.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", response.StatusCode)
	}
	if body != "" {
		t.Errorf("expected an empty body, got %q", body)
	}
}

func TestWriteErrorUsesErrorEnvelope(t *testing.T) {
	response, body := serve(t, func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, http.StatusForbidden, fmt.Errorf("admin access required"))
	})

	if response.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403, got %d", response.StatusCode)
	}
	// types.ErrorResponse in the OpenAPI spec documents exactly this shape.
	if body != "{\"error\":\"admin access required\"}\n" {
		t.Errorf("unexpected body %q", body)
	}
}
