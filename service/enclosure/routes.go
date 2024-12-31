package enclosure

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

type Handler struct {
	store types.EnclosureStore
}

func NewHandler(store types.EnclosureStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/enclosure", h.handleGetEnclosures).Methods(http.MethodGet)
	router.HandleFunc("/enclosure", h.handleCreateEnclosure).Methods(http.MethodPost)
}

func (h *Handler) handleGetEnclosures(w http.ResponseWriter, r *http.Request) {
	enclosureList, err := h.store.GetEnclosures()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, enclosureList)

}
func (h *Handler) handleCreateEnclosure(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var enclosure types.CreateEnclosurePayload
	if err := utils.ParseJSON(r, &enclosure); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(enclosure); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if species exists
	// _, err := h.store.GetSpeciesByComName(species.comName)
	// if err == nil {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("species with common name %s already exists", species.comName))
	// 	return
	// }
	// _, err := h.store.GetSpeciesBySciName(species.sciName)
	// if err == nil {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("species with scientific name %s already exists", species.sciName))
	// 	return
	// }

	// if it doesn't exist, create new species
	err := h.store.CreateEnclosure(types.Enclosure{
		EnclosureName: enclosure.EnclosureName,
		Image:         enclosure.Image,
		Notes:         enclosure.Notes,
		HabitatId:     enclosure.HabitatId,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
