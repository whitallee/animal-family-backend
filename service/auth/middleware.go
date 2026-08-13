package auth

import (
	"fmt"
	"net/http"

	"github.com/whitallee/animal-family-backend/utils"
)

// RequireAdmin rejects the request unless the authenticated caller is an admin.
//
// It must be wrapped by WithJWTAuth, which is what puts the user ID into the
// request context:
//
//	WithJWTAuth(RequireAdmin(handler), userStore)
//
// v1 performed this check inline in each admin handler and, in every handler
// except handleAdminGenerateSpecies, forgot to return after writing the error —
// so non-admins received a 401 and the operation still executed. Gating at the
// middleware layer removes that whole class of bug.
//
// The status is 403 rather than v1's 401: the caller is authenticated (they got
// past WithJWTAuth), they're just not permitted.
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetuserIdFromContext(r.Context())
		if !IsAdmin(userID) {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("admin access required"))
			return
		}

		next(w, r)
	}
}
