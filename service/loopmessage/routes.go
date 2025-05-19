package loopmessage

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/whitallee/animal-family-backend/types"
	"github.com/whitallee/animal-family-backend/utils"
)

type Handler struct {
	store types.LoopMessageStore
}

func NewHandler(store types.LoopMessageStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/loopmessage/inbound", h.handleInboundLoopMessage).Methods(http.MethodPost)
}

func (h *Handler) handleInboundLoopMessage(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var inboundLoopMessagePayload types.InboundLoopMessagePayload
	if err := utils.ParseJSON(r, &inboundLoopMessagePayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload done by other package
	if err := utils.Validate.Struct(inboundLoopMessagePayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// send response before processing
	utils.WriteJSON(w, http.StatusOK, nil)

	// process inbound loop message
	err := h.store.ReceiveLoopMessage(inboundLoopMessagePayload)
	if err != nil {
		// utils.WriteError(w, http.StatusInternalServerError, err)
		fmt.Println(err)
		return
	}

}
