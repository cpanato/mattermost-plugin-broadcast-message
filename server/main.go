package main

import (
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	BotUserKey = "BroadcastBot"
)

func main() {
	plugin.ClientMain(&Plugin{})
}
