package main

import (
	"fmt"
	"log"
	"time"

	"bytes"
	"github.com/dghubble/sling"
	"github.com/xanzy/go-gitlab"
	"strconv"
	"strings"
)

const (
	ChototBaseUrl = "https://gateway.chotot.com"
	AdListingPath = "/v1/public/ad-listing"
	MaxPrice      = 7000000

	// Dingtalk
	AccessToken   = ""

	// Gitlab
	GitlabToken     = ""
	GitlabSnippetID = 1726507
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

func isOldPost(postIdList []string, postId string) bool {
	if len(postIdList) == 0 {
		return false
	}

	for _, id := range postIdList {
		if id == postId {
			return true
		}
	}

	return false
}

func main() {
	var oldPosts []string
	snipString := ""

	git := gitlab.NewClient(nil, GitlabToken)

	for {

		if len(oldPosts) == 0 || snipString == "" {
			// Get latest RAW Snippet
			snippetContent, _, err := git.Snippets.SnippetContent(GitlabSnippetID)
			if err != nil {
				log.Printf("Error getting gitlab Snippet: %s", err)
			}
			snipString := string(snippetContent)

			// Get list of old posts ID
			for _, line := range strings.Split(snipString, "\n") {
				if line == "" {
					continue
				}
				oldPosts = append(oldPosts, strings.Split(line, "|")[0])
			}
		}

		AdListingResult := new(Response)
		_, err := sling.New().Get(ChototBaseUrl + AdListingPath).QueryStruct(Q1Params).ReceiveSuccess(AdListingResult)
		if err != nil {
			log.Printf("Error executing request: %s", err)
		}

		for _, RentalPost := range AdListingResult.Ads {
			if isOldPost(oldPosts, strconv.FormatInt(int64(RentalPost.AdID), 10)) || RentalPost.Price > MaxPrice {
				continue
			}

			Send2DingTalk(
				RentalPost.AdID,
				fmt.Sprintf("https://nha.chotot.com/%d.htm", RentalPost.ListID),
				RentalPost.Subject,
				RentalPost.Price,
				RentalPost.ImageURL,
			)

			var buff bytes.Buffer
			if snipString != "" {
				buff.WriteString(snipString)
				buff.WriteString("\n")
			}
			buff.WriteString(fmt.Sprintf(
				"%d|https://nha.chotot.com/%d.htm|%s|%d",
				RentalPost.AdID,
				RentalPost.ListID,
				RentalPost.Subject,
				RentalPost.Price,
			))
			snipString = buff.String()

			oldPosts = append(oldPosts, fmt.Sprintf("%d", RentalPost.AdID))
		}

		// Update Snippet
		go func() {
			_, _, err = git.Snippets.UpdateSnippet(
				GitlabSnippetID,
				&gitlab.UpdateSnippetOptions{
					Content: &snipString,
				},
			)
			if err != nil {
				log.Printf("Could not update Snippet: %s", err)
			}
		}()

		log.Printf("Finished loop at: %s", time.Now().String())
		time.Sleep(1 * time.Minute)
	}
}
