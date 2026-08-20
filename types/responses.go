package types

import (
	"encoding/json"
	"fmt"
	"time"
)

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

// UserResponse is the wire representation of a user.
//
// User cannot be used directly: its Phone field is a sql.NullString carrying a
// live `json:"phone"` tag and no custom marshaller, so it serialises as
// {"String":"...","Valid":true} rather than a string. v1 shipped that shape
// while the frontend declared `phone: string`, a mismatch nothing caught
// because the types were maintained by hand on both sides.
//
// Password is absent by construction rather than by relying on a `json:"-"`
// tag, so a field added to User later cannot leak into a response by accident.
type UserResponse struct {
	Id        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone" extensions:"x-nullable"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewUserResponse(u *User) UserResponse {
	response := UserResponse{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}

	if u.Phone.Valid {
		phone := u.Phone.String
		response.Phone = &phone
	}

	return response
}

// AuthResponse is returned by login and token refresh.
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// MessageResponse is a bare acknowledgement for endpoints that have nothing
// else to report.
type MessageResponse struct {
	Message string `json:"message" example:"subscription created successfully"`
}

// TaskCompletionResponse reports how many repeating tasks were reset.
type TaskCompletionResponse struct {
	TasksReset int `json:"tasksReset"`
}

// PushSubscriptionResponse describes a registered push subscription.
//
// It deliberately omits the p256dh and auth keys that PushSubscription carries.
// Those exist to encrypt payloads for the endpoint and nothing outside the
// sender needs them, so listing them would widen their exposure for no benefit.
type PushSubscriptionResponse struct {
	SubscriptionId int       `json:"subscriptionId"`
	Endpoint       string    `json:"endpoint"`
	UserAgent      string    `json:"userAgent"`
	CreatedAt      time.Time `json:"createdAt"`
	LastUsed       time.Time `json:"lastUsed"`
}

func NewPushSubscriptionResponse(s *PushSubscription) PushSubscriptionResponse {
	return PushSubscriptionResponse{
		SubscriptionId: s.SubscriptionId,
		Endpoint:       s.Endpoint,
		UserAgent:      s.UserAgent,
		CreatedAt:      s.CreatedAt,
		LastUsed:       s.LastUsed,
	}
}

func NewPushSubscriptionResponses(subscriptions []*PushSubscription) []PushSubscriptionResponse {
	responses := make([]PushSubscriptionResponse, 0, len(subscriptions))
	for _, s := range subscriptions {
		responses = append(responses, NewPushSubscriptionResponse(s))
	}

	return responses
}

// VAPIDKeyResponse carries the public key browsers need to subscribe to push.
type VAPIDKeyResponse struct {
	PublicKey string `json:"publicKey"`
}

// TestNotificationResult is the per-subscription outcome of a test send.
type TestNotificationResult struct {
	SubscriptionId int    `json:"subscriptionId"`
	Endpoint       string `json:"endpoint"`
	HttpStatus     int    `json:"httpStatus"`
	Success        bool   `json:"success"`
	Error          string `json:"error,omitempty"`
}

// TestNotificationResponse reports what happened when a test notification was
// sent to each of the caller's subscriptions.
type TestNotificationResponse struct {
	Message           string                   `json:"message"`
	SubscriptionCount int                      `json:"subscriptionCount"`
	Results           []TestNotificationResult `json:"results"`
}

// AnimalResponse is the wire representation of an animal.
//
// Animal itself cannot be used here. Its LastMessage and MemorialPhotos fields
// are sql.NullString tagged `json:"-"`, and its custom MarshalJSON synthesises
// `lastMessage` and `memorialPhotos` at encode time. The spec generator reflects
// struct tags and cannot see marshalling code, so generating from Animal would
// silently omit both fields — exactly the two the memorial feature depends on.
//
// Declaring the wire shape explicitly also keeps database types out of the API
// contract, and gives nullable fields somewhere to carry `x-nullable` (Swagger
// 2.0 has no nullable keyword; specgen's converter turns the extension into
// OpenAPI 3's `nullable: true`).
//
// Build one with NewAnimalResponse rather than by hand.
type AnimalResponse struct {
	AnimalId        int       `json:"animalId"`
	AnimalName      string    `json:"animalName"`
	SpeciesId       int       `json:"speciesId"`
	EnclosureId     *int      `json:"enclosureId" extensions:"x-nullable"`
	Image           string    `json:"image"`
	Gender          string    `json:"gender"`
	Dob             time.Time `json:"dob"`
	PersonalityDesc string    `json:"personalityDesc"`
	DietDesc        string    `json:"dietDesc"`
	RoutineDesc     string    `json:"routineDesc"`
	ExtraNotes      string    `json:"extraNotes"`
	IsMemorialized  bool      `json:"isMemorialized"`
	LastMessage     *string   `json:"lastMessage" extensions:"x-nullable"`
	MemorialPhotos  []string  `json:"memorialPhotos" extensions:"x-nullable"`
	MemorialDate    time.Time `json:"memorialDate"`
}

// NewAnimalResponse converts a stored animal into its wire representation.
//
// It returns an error only when MemorialPhotos holds text that is not a JSON
// array, which would mean the column was written by something other than the
// memorialise flow.
func NewAnimalResponse(a *Animal) (AnimalResponse, error) {
	response := AnimalResponse{
		AnimalId:        a.AnimalId,
		AnimalName:      a.AnimalName,
		SpeciesId:       a.SpeciesId,
		EnclosureId:     a.EnclosureId,
		Image:           a.Image,
		Gender:          a.Gender,
		Dob:             a.Dob,
		PersonalityDesc: a.PersonalityDesc,
		DietDesc:        a.DietDesc,
		RoutineDesc:     a.RoutineDesc,
		ExtraNotes:      a.ExtraNotes,
		IsMemorialized:  a.IsMemorialized,
		MemorialDate:    a.MemorialDate,
	}

	if a.LastMessage.Valid {
		message := a.LastMessage.String
		response.LastMessage = &message
	}

	if a.MemorialPhotos.Valid {
		var photos []string
		if err := json.Unmarshal([]byte(a.MemorialPhotos.String), &photos); err != nil {
			return AnimalResponse{}, fmt.Errorf("animal %d has malformed memorialPhotos: %w", a.AnimalId, err)
		}
		response.MemorialPhotos = photos
	}

	return response, nil
}

// NewAnimalResponses converts a list of stored animals for a collection route.
func NewAnimalResponses(animals []*Animal) ([]AnimalResponse, error) {
	// Never nil: a nil slice encodes as `null`, and the generated client types
	// list endpoints as returning an array.
	responses := make([]AnimalResponse, 0, len(animals))

	for _, a := range animals {
		response, err := NewAnimalResponse(a)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}

	return responses, nil
}
