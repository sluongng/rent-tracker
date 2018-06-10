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

**URL**: %s

**Subject**: %s

**Price**: %d
`
	content := fmt.Sprintf(contentTemplate, imageURL, adID, url, subject, price)

	err := dingbot.SimpleMarkdownMessage(content).Send(AccessToken)
	if err != nil {
		log.Printf("Something wrong with sending message to dingtalk: %s", err)
	}
}
