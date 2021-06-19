package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/slack-go/slack"
	flag "github.com/spf13/pflag"
)

const (
	shouldEscapeSlackMessage = false
)

type CmdOptions struct {
	TokenFilePath string
	SlackUserID   string
	BotUsername   string
}

func ParseCommandLineOrDie() *CmdOptions {
	opt := &CmdOptions{}

	flag.StringVar(&opt.TokenFilePath, "slack-token-file", "~/.issuemon/slack-token", "read slack bot oauth token from this file")
	flag.StringVar(&opt.SlackUserID, "slack-to", "", "send direct message in slack to this user ID (view profile -> more -> Copy member ID). May also be a channel ID")
	flag.StringVar(&opt.BotUsername, "slack-from", "issuemon", "use this username as sender (requires oauth scope in slack application)")

	flag.Parse()

	if len(opt.SlackUserID) == 0 {
		log.Fatal("slack-to is required")
	}

	return opt
}

func SlackTokenFromFile(opt *CmdOptions) (string, error) {
	path := opt.TokenFilePath
	if strings.HasPrefix(path, "~/") {
		curUser, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("user.Current() error: %v", err)
		}
		path = filepath.Join(curUser.HomeDir, path[2:])
	}

	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("Failed to make absolute path from %q: %v", path, err)
		}
		path = absPath
	}

	slackTokenBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Failed to read gitlab token from file %q: %v", path, err)
	}
	slackToken := strings.TrimSuffix(string(slackTokenBytes), "\n")
	return slackToken, nil
}

func main() {
	opt := ParseCommandLineOrDie()
	slackToken, err := SlackTokenFromFile(opt)
	if err != nil {
		log.Fatalf("slack token read error: %v", err)
	}
	api := slack.New(slackToken)
	//attachment := slack.Attachment{
	//	Pretext: "some pretext",
	//	Text:    "some text",
	//	// Uncomment the following part to send a field too
	//	/*
	//		Fields: []slack.AttachmentField{
	//			slack.AttachmentField{
	//				Title: "a",
	//				Value: "no",
	//			},
	//		},
	//	*/
	//}
	message := fmt.Sprintf("%q: Current time: %v", opt.BotUsername, time.Now())

	channelID, timestamp, err := api.PostMessage(
		opt.SlackUserID,
		slack.MsgOptionText(message, shouldEscapeSlackMessage),
		//slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(false),
		slack.MsgOptionUser(opt.BotUsername),
	)
	if err != nil {
		log.Fatalf("Message send error: %v", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
