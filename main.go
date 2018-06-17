package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dghubble/sling"
	"strconv"
	"strings"
)

const (
	ChototBaseUrl = "https://gateway.chotot.com"
	AdListingPath = "/v1/public/ad-listing"
	MaxPrice      = 7000000
	AccessToken   = "e758dd042bf4c3f825fe46efe817b219c1c3ed4f38df17a6601d384e1ceff6a9"
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

	f, err := os.OpenFile("log/app.log", os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		log.Printf("Error opening log file: %s", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	maxTime := time.Now().Truncate(1 * time.Hour)

	for {
		var oldPosts []string
		for scanner.Scan() {
			lineSlice := strings.Split(scanner.Text(), "|")
			oldPosts = append(oldPosts, lineSlice[0])
		}

		AdListingResult := new(Response)
		_, err := sling.New().Get(ChototBaseUrl + AdListingPath).QueryStruct(Q1Params).ReceiveSuccess(AdListingResult)
		if err != nil {
			log.Printf("Error executing request: %s", err)
		}

		tempMaxTime := maxTime
		for _, RentalPost := range AdListingResult.Ads {
			postTime := time.Unix(RentalPost.ListTime, 0)
			if isOldPost(oldPosts, strconv.FormatInt(int64(RentalPost.AdID), 10)) ||
				RentalPost.Price > MaxPrice ||
				maxTime.After(postTime) ||
				maxTime.Equal(postTime) {
				continue
			}
			if postTime.After(tempMaxTime) {
				tempMaxTime = postTime
			}

			Send2DingTalk(
				RentalPost.AdID,
				fmt.Sprintf("https://nha.chotot.com/%d.htm", RentalPost.ListID),
				RentalPost.Subject,
				RentalPost.Price,
				RentalPost.ImageURL,
			)
			log := fmt.Sprintf("%d|https://nha.chotot.com/%d.htm|%s|%d\n", RentalPost.AdID, RentalPost.ListID, RentalPost.Subject, RentalPost.Price)
			f.WriteString(log)
		}
		maxTime = tempMaxTime

		time.Sleep(1 * time.Minute)
	}
}
