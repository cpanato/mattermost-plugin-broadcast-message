package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

// OnActivate initialize the plugin
func (p *Plugin) OnActivate() error {
	p.API.LogDebug("Activating Broadcast Message plugin")

	if err := p.ensureBotExists(); err != nil {
		return errors.Wrap(err, "failed to ensure bot user exists")
	}

	// Give it a profile picture
	if err := p.API.SetProfileImage(p.botUserID, profileImage); err != nil {
		p.API.LogError("Failed to set profile image for bot", "err", err)
		return errors.Wrap(err, "failed to set profile image for the bot")
	}

	if err := p.API.RegisterCommand(getCommand()); err != nil {
		p.API.LogError("Failed to register the command", "err", err)
		return errors.Wrap(err, "failed to register the command")
	}

	p.API.LogDebug("Broadcast plugin activated")

	return nil
}

func (p *Plugin) ensureBotExists() error {
	// Attempt to find an existing bot
	botUserIDBytes, err := p.API.KVGet(BotUserKey)
	if err != nil {
		return err
	}

	bot := &model.Bot{}
	if botUserIDBytes != nil {
		var appErr *model.AppError
		bot, appErr = p.API.GetBot(string(botUserIDBytes), false)
		if appErr != nil {
			return errors.New(appErr.Error())
		}

		if bot != nil {
			p.botUserID = bot.UserId
			return nil
		}
	}

	if botUserIDBytes == nil || bot == nil {
		// Create a bot since one doesn't exist
		p.API.LogDebug("Creating bot for broadcast plugin")

		bot, err := p.API.CreateBot(&model.Bot{
			Username:    "broadcast_bot",
			DisplayName: "Broadcast Bot",
			Description: "Created by the Broadcast plugin.",
		})
		if err != nil {
			return err
		}

		p.API.LogDebug("Bot created for broadcast plugin")

		// Save the bot ID
		err = p.API.KVSet(BotUserKey, []byte(bot.UserId))
		if err != nil {
			return err
		}
		p.botUserID = bot.UserId
	}

	return nil
}
