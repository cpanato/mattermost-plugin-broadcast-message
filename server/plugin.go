package main

import (
	"net/http"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	botUserID string
}

// ServeHTTP allows the plugin to implement the http.Handler interface. Requests destined for the
// /plugins/{id} path will be routed to the plugin.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" || strings.Compare(token, p.configuration.Token) != 0 {
		errorMessage := "Invalid or missing token"
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	switch r.URL.Path {
	case "/api/dialog":
		p.handleDialog(w, r)
	default:
		http.NotFound(w, r)
	}
}
