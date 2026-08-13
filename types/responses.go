package types

// Response types used by the v2 API.
//
// These exist so the OpenAPI generator has concrete types to reference. v1
// wrote several responses as ad hoc map[string]interface{} literals, which
// produce correct JSON but leave nothing for the spec (and therefore the
// generated TypeScript client) to describe.

// ErrorResponse is the body written by utils.WriteError for every non-2xx
// response. Its shape must stay in sync with that function.
type ErrorResponse struct {
	Error string `json:"error" example:"admin access required"`
}
