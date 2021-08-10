// This file is automatically generated. Do not modify it manually.

package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

var manifest *model.Manifest

const manifestStr = `
{
  "id": "com.cpanato.broadcast",
  "name": "BroadCast Message",
  "description": "This plugin broadcast a message to all users in the Mattermost instance.",
  "version": "0.1.1",
  "min_server_version": "5.12.0",
  "server": {
    "executables": {
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "linux-amd64": "server/dist/plugin-linux-amd64",
      "windows-amd64": "server/dist/plugin-windows-amd64.exe"
    },
    "executable": ""
  },
  "settings_schema": {
    "header": "",
    "footer": "",
    "settings": [
      {
        "key": "Token",
        "display_name": "Token:",
        "type": "generated",
        "help_text": "The token used validate the requests.",
        "placeholder": "",
        "default": null
      }
    ]
  }
}
`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}
