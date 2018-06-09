package main

import (
	"fmt"
	"log"

	"github.com/sluongng/dingbot"
)

func Send2DingTalk(adID int32, url string, subject string, price int64, imageURL string) {

	contentTemplate := `

![](%s)

**AdID**: %d

**URL**: [LINK](%s)

**Subject**: %s

**Price**: %d
`
	content := fmt.Sprintf(contentTemplate, imageURL, adID, url, subject, price)

	ChatBotService := dingbot.NewClient(AccessToken).RobotService
	mdMessage := &dingbot.MarkdownMessage{
		MsgType: "markdown",
		Markdown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{
			Title: "New Posting",
			Text:  content,
		},
	}

	err := ChatBotService.SendMarkdown(mdMessage)
	if err != nil {
		log.Printf("Something wrong with sending message to dingtalk: %s", err)
	}
}
