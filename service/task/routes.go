package task

import (
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/types"
)

type Handler struct {
	store          types.TaskStore
	userStore      types.UserStore
	animalStore    types.AnimalStore
	enclosureStore types.EnclosureStore
}

func NewHandler(store types.TaskStore, userStore types.UserStore, animalStore types.AnimalStore, enclosureStore types.EnclosureStore) *Handler {
	return &Handler{store: store, userStore: userStore, animalStore: animalStore, enclosureStore: enclosureStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// user routes

	// admin routes

}
