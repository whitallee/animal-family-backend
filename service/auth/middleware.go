package auth

import (
	"context"
	"fmt"
	"log"
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

// ResourceIDKey holds the validated path-parameter ID that RequireOwnership
// checked, so handlers do not parse it a second time.
const ResourceIDKey contextKey = "resourceID"

// OwnershipChecker reports whether userID owns resourceID. A false with a nil
// error means "not owned"; a non-nil error means the check itself failed.
type OwnershipChecker func(resourceID int, userID int) (bool, error)

// RequireOwnership parses the named path parameter, verifies the authenticated
// caller owns that resource, and passes the validated ID to the handler through
// the request context.
//
// It must be wrapped by WithJWTAuth, which supplies the user ID:
//
//	WithJWTAuth(RequireOwnership("id", store.UserOwnsAnimal, handler), userStore)
//
// v1 repeated this check inline in every handler (six times in
// service/animal/routes.go alone) and answered 400 for an ownership failure,
// which reads as "your request was malformed" rather than "not yours". It also
// could not tell a permission failure from a database error, because the store
// returned a plain error for both.
//
// A resource that does not exist and one owned by somebody else both yield 403,
// so the response does not reveal which.
func RequireOwnership(param string, check OwnershipChecker, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.ParseIDParam(r, param)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		userID := GetuserIdFromContext(r.Context())

		owned, err := check(id, userID)
		if err != nil {
			log.Printf("ownership check failed (resource %d, user %d): %v", id, userID, err)
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("could not verify access to this resource"))
			return
		}

		if !owned {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("you do not have access to this resource"))
			return
		}

		next(w, r.WithContext(context.WithValue(r.Context(), ResourceIDKey, id)))
	}
}

// ResourceIDFromContext returns the ID validated by RequireOwnership. It
// returns -1 when the handler was not wrapped, matching GetuserIdFromContext.
func ResourceIDFromContext(ctx context.Context) int {
	id, ok := ctx.Value(ResourceIDKey).(int)
	if !ok {
		return -1
	}

	return id
}
