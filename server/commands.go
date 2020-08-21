package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	broadcast = "/broadcast-message"
)

func getCommand() *model.Command {
	return &model.Command{
		Trigger:      "broadcast-message",
		DisplayName:  "Broadcast Message",
		Description:  "Broadcast a message to all users. Only sys-admins can use this command.",
		AutoComplete: true,
	}
}

// ExecuteCommand execute the slash command
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	hasPermission := p.API.HasPermissionTo(args.UserId, model.PERMISSION_MANAGE_SYSTEM)

	if !hasPermission {
		permissionErrPost := &model.Post{
			UserId:    p.botUserID,
			ChannelId: args.ChannelId,
			Message:   "Only System Administrators have permission to use this command.",
		}

		_ = p.API.SendEphemeralPost(args.UserId, permissionErrPost)

		return &model.CommandResponse{}, nil
	}

	split := strings.Fields(args.Command)
	command := split[0]
	if command != broadcast {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Not a valid command"), nil
	}

	user, uErr := p.API.GetUser(args.UserId)
	if uErr != nil {
		return &model.CommandResponse{}, uErr
	}
	if strings.Trim(args.Command, " ") == broadcast {
		p.InteractiveSchedule(args.TriggerId, user)
		return &model.CommandResponse{}, nil
	}

	message := strings.Join(split[1:], " ")

	err := p.broadcastMessage(args.TeamId, args.ChannelId, args.UserId, message)
	if err != nil {
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Error sending the broadcast message"), nil
	}

	return &model.CommandResponse{}, nil
}

func (p *Plugin) createBotDMPost(userID string, post *model.Post) (*model.Post, *model.AppError) {
	channel, err := p.API.GetDirectChannel(userID, p.botUserID)
	if err != nil {
		p.API.LogError("Couldn't get bot's DM channel", "user_id", userID, "err", err)
		return nil, err
	}

	post.UserId = p.botUserID
	post.ChannelId = channel.Id

	created, err := p.API.CreatePost(post)
	if err != nil {
		p.API.LogError("Couldn't send bot DM", "user_id", userID, "err", err)
		return nil, err
	}

	return created, nil
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Type:         model.POST_DEFAULT,
	}
}

func (p *Plugin) broadcastMessage(teamID, channelID, userID, message string) error {
	broadcastMessage := &model.Post{
		UserId:  p.botUserID,
		Message: message,
	}

	page := 0
	for {
		users, appErr := p.API.GetUsersInTeam(teamID, page, 200)
		if appErr != nil {
			return errors.New(appErr.Error())
		}

		if users == nil {
			break
		}

		for _, user := range users {
			if !user.IsBot {
				_, appErr := p.createBotDMPost(user.Id, broadcastMessage)
				if appErr != nil {
					return errors.New(appErr.Error())
				}
			}
		}

		page++
	}

	OkPost := &model.Post{
		UserId:    p.botUserID,
		ChannelId: channelID,
		Message:   fmt.Sprintf("Message broadcasted - message %s", message),
	}

	_ = p.API.SendEphemeralPost(userID, OkPost)

	return nil
}

func (p *Plugin) InteractiveSchedule(triggerID string, user *model.User) {
	config := p.API.GetConfig()
	siteURLPort := *config.ServiceSettings.ListenAddress
	dialogRequest := model.OpenDialogRequest{
		TriggerId: triggerID,
		URL:       fmt.Sprintf("http://localhost%v/plugins/%v/api/dialog?token=%s", siteURLPort, manifest.Id, p.configuration.Token),
		Dialog: model.Dialog{
			Title:       "Broadcast an important message",
			CallbackId:  model.NewId(),
			SubmitLabel: "Broadcast",
			Elements: []model.DialogElement{
				{
					DisplayName: "Message to broadcast",
					Name:        "message",
					Type:        "textarea",
					HelpText:    "Message to broadcast",
					Optional:    false,
				},
			},
		},
	}
	if pErr := p.API.OpenInteractiveDialog(dialogRequest); pErr != nil {
		p.API.LogError("Failed opening interactive dialog " + pErr.Error())
	}
}
