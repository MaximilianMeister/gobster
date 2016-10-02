// Package main
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/thoj/go-ircevent"

	"github.com/MaximilianMeister/gobster/api"
	"github.com/MaximilianMeister/gobster/feed"
)

// Configuration holds static config data
// It needs to look like:
// {
//   "irc_server": "irc.freenode.net",
//   "irc_port": 6666,
//   "irc_channel": "#gobster",
//   "irc_botnick": "gobster",
//   "irc_botname": "Gobbi",
//   "irc_welcome_msg": "Hello",
//   "bot_set_msg": "Successfully added message",
//   "bot_set_error_msg": "Could not add message",
//   "bot_get_error_msg": "Could not find message",
//   "default_bucket": "bucket" # the default key/bucket from where to fetch a random entry
// }
type Configuration struct {
	IrcServer      string `json:"irc_server"`
	IrcPort        string `json:"irc_port"`
	IrcChannel     string `json:"irc_channel"`
	IrcBotnick     string `json:"irc_botnick"`
	IrcBotname     string `json:"irc_botname"`
	IrcWelcomeMsg  string `json:"irc_welcome_msg"`
	BotSetMsg      string `json:"bot_set_msg"`
	BotSetErrorMsg string `json:"bot_set_error_msg"`
	BotGetErrorMsg string `json:"bot_get_error_msg"`
	DefaultBucket  string `json:"default_bucket"`
}

// GetConfig returns a type Configuration with values defined in gobster.json
func GetConfig() (Configuration, error) {
	config := &Configuration{}
	path, err := filepath.Abs("./gobster.json")
	if err != nil {
		fmt.Printf("File path error: %v\n", err)
		os.Exit(1)
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	if err = json.Unmarshal(file, config); err != nil {
		return Configuration{}, err
	}

	return *config, nil
}

func main() {
	config, err := GetConfig()
	if err != nil {
		fmt.Println("Cannot read configuration file: ", err)
		os.Exit(1)
	}

	ircobj := irc.IRC(config.IrcBotnick, config.IrcBotname)
	ircobj.Connect(fmt.Sprintf("%s:%s", config.IrcServer, config.IrcPort))
	ircobj.Join(config.IrcChannel)
	ircobj.Privmsg(config.IrcChannel, config.IrcWelcomeMsg)

	ircobj.AddCallback("PRIVMSG", func(event *irc.Event) {
		if !strings.Contains(event.Message(), "!") {
			if strings.Contains(event.Message(), config.IrcBotnick) {
				quote, err := feed.Get(config.DefaultBucket)
				if err != nil {
					// no default bucket exists, so give a generic answer
					ircobj.Privmsg(config.IrcChannel, "There is nothing more to say")
					return
				}
				ircobj.Privmsg(config.IrcChannel, fmt.Sprintf("%s: %s", event.Nick, quote))
				return
			} else {
				return
			}
		}

		messages := strings.Split(event.Message(), " ")
		command := strings.Split(messages[0], "!")[1]
		argc := len(messages) - 1
		recipient := messages[argc]

		if (len(messages) >= 2) && (messages[1] == "add") {
			newQuote := strings.Join(messages[2:len(messages)][:], " ")
			err := feed.Set(command, newQuote)
			if err != nil {
				ircobj.Privmsg(config.IrcChannel, config.BotSetErrorMsg)
			} else {
				ircobj.Privmsg(config.IrcChannel, config.BotSetMsg)
			}
		} else {
			quote, err := feed.Get(command)
			if err != nil {
				quote = config.BotGetErrorMsg
			}

			if argc == 0 {
				ircobj.Privmsg(config.IrcChannel, quote)
			} else {
				ircobj.Privmsg(config.IrcChannel, fmt.Sprintf("%s: %s", recipient, quote))
			}
		}
	})

	api.Serve()
}
