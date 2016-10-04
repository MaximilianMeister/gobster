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
		if event.Message() == "!help" {
			ircobj.Privmsg(config.IrcChannel, "Gobster Usage:")
			ircobj.Privmsg(config.IrcChannel, "Available commands:")
			ircobj.Privmsg(config.IrcChannel, "  add: Add a message to a [sub]command")
			ircobj.Privmsg(config.IrcChannel, "    '!command add This is a quote'")
			ircobj.Privmsg(config.IrcChannel, "    '!command subcommand add This is a quote'")
			return
		}
		if strings.Contains(event.Message(), config.IrcBotnick) { // someone talks to the bot directly
			quote, err := feed.Get(config.DefaultBucket)
			if err != nil {
				// no default bucket exists, so give a generic answer
				ircobj.Privmsg(config.IrcChannel, "There is nothing more to say")
				return
			}
			ircobj.Privmsg(config.IrcChannel, fmt.Sprintf("%s: %s", event.Nick, quote))
			return
		}

		messages := strings.Split(event.Message(), " ")

		if !strings.HasPrefix(event.Message(), "!") { // there is no ! prefix
			return
		}

		command := strings.Split(messages[0], "!")[1]

		if command == "" { // nothing to do
			return
		}

		subcommand := ""
		argc := len(messages) - 1
		recipient := messages[argc]

		// determin if there is a subcommand
		if len(messages) >= 2 {
			subcommand = fmt.Sprintf("%s_%s", command, messages[1])
			// make sure the subcommand exists
			quote, err := feed.Get(subcommand)
			if (quote == "") || (err != nil) {
				subcommand = ""
			}
		}

		// determine if we want to add a string to a main or a subcommand
		newQuote := ""
		if len(messages) >= 2 {
			if messages[1] == "add" {
				newQuote = strings.Join(messages[2:len(messages)][:], " ")
			} else if (len(messages) >= 3) && (messages[2] == "add") {
				newQuote = strings.Join(messages[3:len(messages)][:], " ")
			}
		}

		if newQuote != "" {
			var err error
			if subcommand != "" {
				err = feed.Set(subcommand, newQuote)
			} else {
				err = feed.Set(command, newQuote)
			}
			if err != nil {
				ircobj.Privmsg(config.IrcChannel, config.BotSetErrorMsg)
				return
			} else { // successfully set
				ircobj.Privmsg(config.IrcChannel, config.BotSetMsg)
				return
			}
		} else if (len(messages) >= 2) && subcommand != "" { // get string from a subcommand
			quote, err := feed.Get(subcommand)
			if err != nil {
				quote = config.BotGetErrorMsg
			}
			if argc >= 2 { // there is a recipient
				ircobj.Privmsg(config.IrcChannel, fmt.Sprintf("%s: %s", recipient, quote))
			} else { // there is no recipient
				ircobj.Privmsg(config.IrcChannel, quote)
			}
		} else { // get string from a main command
			quote, err := feed.Get(command)
			if err != nil {
				quote = config.BotGetErrorMsg
			}

			if argc == 0 { // there is no recipient
				ircobj.Privmsg(config.IrcChannel, quote)
			} else { // there is a recipient
				ircobj.Privmsg(config.IrcChannel, fmt.Sprintf("%s: %s", recipient, quote))
			}
		}
	})

	// start the api server
	api.Serve()
}
