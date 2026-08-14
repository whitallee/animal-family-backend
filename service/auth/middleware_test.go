package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
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

// ownedRequest routes through mux so mux.Vars is populated as in production,
// and injects the user ID the way WithJWTAuth does.
func runOwnership(t *testing.T, url string, userID int, check OwnershipChecker) (*httptest.ResponseRecorder, bool, int) {
	t.Helper()

	var (
		handlerRan bool
		seenID     int
	)

	router := mux.NewRouter()
	router.HandleFunc("/animals/{id}", RequireOwnership("id", check, func(w http.ResponseWriter, r *http.Request) {
		handlerRan = true
		seenID = ResourceIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, url, nil)
	request = request.WithContext(context.WithValue(request.Context(), UserKey, userID))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder, handlerRan, seenID
}

func ownedBy(owner int) OwnershipChecker {
	return func(resourceID int, userID int) (bool, error) {
		return userID == owner, nil
	}
}

func TestRequireOwnershipAllowsOwner(t *testing.T) {
	recorder, handlerRan, seenID := runOwnership(t, "/animals/42", 7, ownedBy(7))

	if !handlerRan {
		t.Error("expected the handler to run for the owner")
	}
	if recorder.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", recorder.Code)
	}
	// The handler must be able to read the ID without re-parsing it.
	if seenID != 42 {
		t.Errorf("expected the validated id 42 in context, got %d", seenID)
	}
}

func TestRequireOwnershipBlocksNonOwner(t *testing.T) {
	recorder, handlerRan, _ := runOwnership(t, "/animals/42", 7, ownedBy(999))

	if handlerRan {
		t.Error("handler ran for a non-owner; the request must be stopped")
	}
	if recorder.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", recorder.Code)
	}
}

// A failed lookup must not read as a permission decision: answering 403 would
// tell the caller the resource is not theirs when the truth is the database
// could not be reached. v1 could not tell these apart at all.
func TestRequireOwnershipReportsLookupFailureAsServerError(t *testing.T) {
	failing := func(resourceID int, userID int) (bool, error) {
		return false, fmt.Errorf("connection refused")
	}

	recorder, handlerRan, _ := runOwnership(t, "/animals/42", 7, failing)

	if handlerRan {
		t.Error("handler ran despite the ownership check failing")
	}
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", recorder.Code)
	}
}

func TestRequireOwnershipRejectsInvalidID(t *testing.T) {
	checkCalls := 0
	counting := func(resourceID int, userID int) (bool, error) {
		checkCalls++
		return true, nil
	}

	recorder, handlerRan, _ := runOwnership(t, "/animals/0", 7, counting)

	if handlerRan {
		t.Error("handler ran for an invalid id")
	}
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", recorder.Code)
	}
	// The ID is validated before any lookup, so a bad ID never reaches the store.
	if checkCalls != 0 {
		t.Errorf("expected no ownership lookups for an invalid id, got %d", checkCalls)
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
