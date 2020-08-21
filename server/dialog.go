package main

import (
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (p *Plugin) handleDialog(w http.ResponseWriter, r *http.Request) {
	p.API.LogInfo("Received dialog action")

	request := model.SubmitDialogRequestFromJson(r.Body)

	message := request.Submission["message"]
	go func() {
		_ = p.broadcastMessage(request.TeamId, request.ChannelId, request.UserId, message.(string))
	}()
}
