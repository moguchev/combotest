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
func SetUsersHandler(router *mux.Router) {
	handler := &UsersHandler{}

	router.HandleFunc("/users", handler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc(fmt.Sprintf("/users/{%s}", userIDParam), handler.ApproveUser).Methods(http.MethodPost)
	router.HandleFunc("/auth", handler.Authorize).Methods(http.MethodPost)
}

func (h *UsersHandler) Authorize(w http.ResponseWriter, r *http.Request) {

}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

}

func (h *UsersHandler) ApproveUser(w http.ResponseWriter, r *http.Request) {

}
