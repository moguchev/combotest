package delivery

import (
	"combotest/internal/interfaces/events"
	"combotest/internal/models"
	"combotest/pkg/utils"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	eventIDParam = "id"
)

// EventsHandler  represent the http handler for events
type EventsHandler struct {
	EventsUC events.Usecase
	Log      *logrus.Logger
}

// SetEventsHandler will set handlers
func (handler *EventsHandler) SetEventsHandler(router *mux.Router) {
	router.HandleFunc("/", handler.GetEvents).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/{%s}", eventIDParam), handler.UpdateEvent).Methods(http.MethodPatch)
	router.HandleFunc("/incedent", handler.SetIncedentInEvents).Methods(http.MethodPost)
}

const (
	limitParam  = "limit"
	offsetParam = "offset"

	systemNameParam      = "system_name"
	incidentParam        = "incident"
	createdAtAfterParam  = "created_at_after"
	createdAtBeforeParam = "created_at_before"
	eventIDQueryParam    = "event_id"
)

func getEventsFilter(values url.Values) models.EventsFilter {
	f := models.EventsFilter{}

	// CreatedAt * TimeInterval // BETWEEN

	l, ok := values[limitParam]
	if ok {
		lim, err := strconv.ParseUint(l[0], 10, 64)
		if err == nil {
			lim32 := uint32(lim)
			f.Limit = &lim32
		}
	}

	o, ok := values[offsetParam]
	if ok {
		off, err := strconv.ParseUint(o[0], 10, 64)
		if err == nil {
			off32 := uint32(off)
			f.Offset = &off32
		}
	}

	eid, ok := values[eventIDQueryParam]
	if ok {
		i, err := strconv.ParseInt(eid[0], 10, 64)
		if err == nil {
			id := int(i)
			f.EventID = &id
		}
	}

	inc, ok := values[incidentParam]
	if ok {
		incedent, err := strconv.ParseBool(inc[0])
		if err == nil {
			f.Incident = &incedent
		}
	}

	after, ok := values[createdAtAfterParam]
	if ok {
		timeAfter, err := time.Parse(time.RFC3339, after[0])
		if err == nil {
			if f.CreatedAt != nil {
				f.CreatedAt.Begin = &timeAfter
			} else {
				f.CreatedAt = &models.TimeInterval{Begin: &timeAfter}
			}
		}
	}

	before, ok := values[createdAtAfterParam]
	if ok {
		timeBefore, err := time.Parse(time.RFC3339, before[0])
		if err == nil {
			if f.CreatedAt != nil {
				f.CreatedAt.Begin = &timeBefore
			} else {
				f.CreatedAt = &models.TimeInterval{Begin: &timeBefore}
			}
		}
	}

	return f
}

func (h *EventsHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	log := h.Log.WithField("handler", "GetEvents")
	ctx := r.Context()

	filter := getEventsFilter(r.URL.Query())

	total, events, err := h.EventsUC.GetEvents(ctx, filter)
	if err != nil {
		log.WithError(err).Error("get users: %w", err)
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	type Response struct {
		Total  uint32         `json:"total"`
		Events []models.Event `json:"events"`
	}

	utils.RespondWithJSON(w, http.StatusOK, Response{Total: total, Events: events})
}

func (h *EventsHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {

}

func (h *EventsHandler) SetIncedentInEvents(w http.ResponseWriter, r *http.Request) {

}
