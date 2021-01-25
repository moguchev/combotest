package delivery

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	eventIDParam = "id"
)

// EventsHandler  represent the http handler for events
type EventsHandler struct {
}

// SetEventsHandler will set handlers
func (handler *EventsHandler) SetEventsHandler(router *mux.Router) {
	router.HandleFunc("/", handler.GetEvents).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/{%s}", eventIDParam), handler.UpdateEvent).Methods(http.MethodPatch)
	router.HandleFunc("/incedent", handler.SetIncedentInEvents).Methods(http.MethodPost)
}

func (h *EventsHandler) GetEvents(w http.ResponseWriter, r *http.Request) {

}

func (h *EventsHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {

}

func (h *EventsHandler) SetIncedentInEvents(w http.ResponseWriter, r *http.Request) {

}
