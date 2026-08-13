package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// adminUserID must match an entry in the admins slice in admin.go.
const (
	adminUserID    = 6
	nonAdminUserID = 99999
)

func requestWithUser(userID int) *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/habitats", nil)
	return r.WithContext(context.WithValue(r.Context(), UserKey, userID))
}

func TestRequireAdminAllowsAdmin(t *testing.T) {
	called := false
	handler := RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	rr := httptest.NewRecorder()
	handler(rr, requestWithUser(adminUserID))

	if !called {
		t.Error("expected wrapped handler to run for an admin")
	}
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}

// TestRequireAdminBlocksNonAdmin covers the bug this middleware exists to
// prevent: v1 wrote the error but did not return, so the handler still ran and
// the operation succeeded for non-admins. Asserting the handler never runs is
// the point of this test - a status-only assertion would have passed against
// the v1 code too.
func TestRequireAdminBlocksNonAdmin(t *testing.T) {
	called := false
	handler := RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	rr := httptest.NewRecorder()
	handler(rr, requestWithUser(nonAdminUserID))

	if called {
		t.Error("wrapped handler ran for a non-admin; the request must be stopped")
	}
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}
}

// A request that never passed through WithJWTAuth has no user ID in context,
// which GetuserIdFromContext reports as -1. That must not be treated as admin.
func TestRequireAdminBlocksMissingUserContext(t *testing.T) {
	called := false
	handler := RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	rr := httptest.NewRecorder()
	handler(rr, httptest.NewRequest(http.MethodPost, "/habitats", nil))

	if called {
		t.Error("wrapped handler ran without an authenticated user")
	}
	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}
}
