package controllers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/sniperkit/watchub/config"
	"github.com/sniperkit/watchub/shared/pages"
)

// Contact ctrl
type Contact struct {
	Base
}

// NewContact ctrl
func NewContact(
	config config.Config,
	session sessions.Store,
) *Contact {
	return &Contact{
		Base: Base{
			config:  config,
			session: session,
		},
	}
}

// Handler handles /Contact
func (ctrl *Contact) Handler(w http.ResponseWriter, r *http.Request) {
	pages.Render(w, "contact", ctrl.sessionData(w, r))
}
