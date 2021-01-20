package delivery

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	userIDParam = "user_id"
)

// UsersHandler  represent the http handler for events
type UsersHandler struct {
}

// SetUsersHandler will set handlers
func (handler *UsersHandler) SetUsersHandler(router *mux.Router) {
	router.HandleFunc("/users", handler.GetUsers).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/users/{%s}", userIDParam), handler.ApproveUser).Methods(http.MethodPost)
}

func (h *UsersHandler) ApproveUser(w http.ResponseWriter, r *http.Request) {

}

func (h *UsersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
}
