package delivery

import (
	"combotest/internal/interfaces/users"
	"combotest/internal/models"
	"combotest/pkg/utils"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	userIDParam = "user_id"
)

// UsersHandler  represent the http handler for events
type UsersHandler struct {
	UsersUC users.Usecase
	Log     *logrus.Logger
}

// SetUsersHandler will set handlers
func (handler *UsersHandler) SetUsersHandler(router *mux.Router) {
	router.HandleFunc("/", handler.GetUsers).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/{%s}", userIDParam), handler.ApproveUser).Methods(http.MethodPatch)
}

func (h *UsersHandler) ApproveUser(w http.ResponseWriter, r *http.Request) {
	log := h.Log.WithField("handler", "ApproveUser")
	ctx := r.Context()

	userID := mux.Vars(r)[userIDParam]
	if err := h.UsersUC.ApproveUser(ctx, userID); err != nil {
		log.WithError(err).Error("approve user: %w", err)
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UsersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	log := h.Log.WithField("handler", "GetUsers")
	ctx := r.Context()

	f := models.UserFilter{}

	// TODO parse URL Query params

	total, users, err := h.UsersUC.GetUsers(ctx, f)
	if err != nil {
		log.WithError(err).Error("get users: %w", err)
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	type Response struct {
		Total uint32        `json:"total"`
		Users []models.User `json:"users"`
	}

	utils.RespondWithJSON(w, http.StatusOK, Response{Total: total, Users: users})
}
