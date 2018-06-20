package controllers

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/sniperkit/watchub/config"
	"github.com/sniperkit/watchub/datastore"
	"github.com/sniperkit/watchub/shared/pages"
)

// Schedule ctrl
type Schedule struct {
	Base
	store datastore.Datastore
}

// NewSchedule ctrl
func NewSchedule(
	config config.Config,
	session sessions.Store,
	store datastore.Datastore,
) *Schedule {
	return &Schedule{
		Base: Base{
			config:  config,
			session: session,
		},
		store: store,
	}
}

// Handler handles /Schedule
func (ctrl *Schedule) Handler(w http.ResponseWriter, r *http.Request) {
	session, _ := ctrl.session.Get(r, ctrl.config.SessionName)
	id, _ := session.Values["user_id"].(int)
	if session.IsNew || id == 0 {
		http.Error(w, "not logged in", http.StatusForbidden)
		return
	}
	if err := ctrl.store.Schedule(int64(id), time.Now()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pages.Render(w, "scheduled", ctrl.sessionData(w, r))
}
