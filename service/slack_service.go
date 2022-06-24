package service

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func InitSlackClient() *slack.Client {
	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
	return client
}

func InitSocketClient(client *slack.Client) *socketmode.Client {
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	return socketClient
}

func SendSlackMessage(client *slack.Client, text string) {
	channelID := os.Getenv("SLACK_CHANNEL_ID")

	attachment := slack.Attachment{
		Pretext: text,
	}
	_, timestamp, err := client.PostMessage(
		channelID,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Message sent at %s", timestamp)
}

func UploadFile(client *slack.Client, fileName string, location string) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	executionPath := filepath.Dir(ex)

	file, err := os.Open(path.Join(path.Dir(executionPath)+"/buzzwordBot/", location))

	channelID := os.Getenv("SLACK_CHANNEL_ID")

	if err != nil {
		fmt.Printf("opening file error: %+v", err)
	}

	params := slack.FileUploadParameters{
		Reader:   file,
		Channels: []string{channelID},
		Filename: fileName,
		Title:    "",
	}

	_, err = client.UploadFile(params)

	if err != nil {
		fmt.Printf("uploading file error: %+v", err)
	}
}
