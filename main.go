package main

import (
	"github.com/dghubble/sling"
	"log"
	"fmt"
	"time"
)

const (
	ChototBaseUrl = "https://gateway.chotot.com"
	AdListingPath = "/v1/public/ad-listing"
	MaxPrice = 7000000

)

var (
	Q1Params = &Params{
		Region: 13,
		Area:96,
		Cg:1010,
		Page:1,
		Limit:20,
		O: 40,
	}
)

type Params struct {
	Region int `url:"region,omitempty"`
	Area int `url:"area,omitempty"`
	Cg int `url:"cg,omitempty"`
	Page int `url:"page,omitempty"`
	Limit int `url:"limit,omitempty"`
	O int `url:"o,omitempty"`
}

type AdPost struct {
	AdID int32 `json:"ad_id"`
	ListID int32 `json:"list_id"`
	ListTime int64 `json:"list_time"`
	Subject string `json:"subject"`
	Price int64 `json:"price"`
	IsCompanyAd bool `json:"company_ad"`
}

type Response struct {
	Total int `json:"total"`
	Ads []AdPost `json:"ads"`
}

func main() {

	maxTime := time.Now().Truncate(1*time.Hour)

	for {
		//log.Printf("MaxTime is %d", maxTime)

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

			// TODO: Launch a chat webhook here.
		}
		maxTime = tempMaxTime
		//fmt.Printf("New MaxTime is %d", maxTime)

		time.Sleep(10*time.Second)
	}
}
