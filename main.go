package main

import (
	"fmt"
	"github.com/dghubble/sling"
	"github.com/sluongng/dingbot"
	"log"
	"time"
)

const (
	ChototBaseUrl = "https://gateway.chotot.com"
	AdListingPath = "/v1/public/ad-listing"
	MaxPrice      = 7000000
	AccessToken   = "123456"
)

var (
	Q1Params = &Params{
		Region: 13,
		Area:   96,
		Cg:     1010,
		Page:   1,
		Limit:  20,
		O:      40,
	}
)

type Params struct {
	Region int `url:"region,omitempty"`
	Area   int `url:"area,omitempty"`
	Cg     int `url:"cg,omitempty"`
	Page   int `url:"page,omitempty"`
	Limit  int `url:"limit,omitempty"`
	O      int `url:"o,omitempty"`
}

type AdPost struct {
	AdID        int32  `json:"ad_id"`
	ListID      int32  `json:"list_id"`
	ListTime    int64  `json:"list_time"`
	Subject     string `json:"subject"`
	Price       int64  `json:"price"`
	IsCompanyAd bool   `json:"company_ad"`
	ImageURL    string `json:"image"`
}

type Response struct {
	Total int      `json:"total"`
	Ads   []AdPost `json:"ads"`
}

func main() {

	maxTime := time.Now().Truncate(1 * time.Hour)

	for {
		AdListingResult := new(Response)
		_, err := sling.New().Get(ChototBaseUrl + AdListingPath).QueryStruct(Q1Params).ReceiveSuccess(AdListingResult)
		if err != nil {
			log.Fatalf("Error executing request: %s", err)
		}

		tempMaxTime := maxTime
		for _, RentalPost := range AdListingResult.Ads {
			postTime := time.Unix(RentalPost.ListTime, 0)
			if RentalPost.Price > MaxPrice || maxTime.After(postTime) || maxTime.Equal(postTime) {
				continue
			}
			if postTime.After(tempMaxTime) {
				tempMaxTime = postTime
			}

			fmt.Printf("AdID: %d\n ", RentalPost.AdID)
			fmt.Printf("URL: https://nha.chotot.com/%d.htm \n", RentalPost.ListID)
			fmt.Printf("Subject: %s \n", RentalPost.Subject)
			fmt.Printf("Price: %d \n", RentalPost.Price)
			fmt.Printf("*** \n")
			Send2DingTalk(
				RentalPost.AdID,
				fmt.Sprintf("https://nha.chotot.com/%d.htm", RentalPost.ListID),
				RentalPost.Subject,
				RentalPost.Price,
				RentalPost.ImageURL,
			)

			// TODO: Launch a chat webhook here.
		}
		maxTime = tempMaxTime

		time.Sleep(1 * time.Minute)
	}
}

func Send2DingTalk(adID int32, url string, subject string, price int64, imageURL string) {
	contentTemplate := `

![](%s)

**AdID**: %d

**URL**: %s

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
