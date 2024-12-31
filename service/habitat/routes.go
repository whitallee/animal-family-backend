package habitat

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

type Handler struct {
	store types.HabitatStore
}

func NewHandler(store types.HabitatStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/habitat", h.handleGetHabitats).Methods(http.MethodGet)
	router.HandleFunc("/habitat", h.handleCreateHabitat).Methods(http.MethodPost)
}

func (h *Handler) handleGetHabitats(w http.ResponseWriter, r *http.Request) {
	habitatsList, err := h.store.GetHabitats()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, habitatsList)

}
func (h *Handler) handleCreateHabitat(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var habitat types.CreateHabitatPayload
	if err := utils.ParseJSON(r, &habitat); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(habitat); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if habitat exists
	// _, err := h.store.GetHabitatByName(habitat.habitatName)
	// if err == nil {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("habitat with name %s already exists", habitat.habitatName))
	// 	return
	// }

	// if it doesn't exist, create new habitat
	err := h.store.CreateHabitat(types.Habitat{
		HabitatName:    habitat.HabitatName,
		HabitatDesc:    habitat.HabitatDesc,
		Image:          habitat.Image,
		Humidity:       habitat.Humidity,
		DayTempRange:   habitat.DayTempRange,
		NightTempRange: habitat.NightTempRange,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
