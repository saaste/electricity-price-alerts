package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	arg "github.com/alexflint/go-arg"
)

type PriceData struct {
	Price     float32
	StartDate time.Time
	EndDate   time.Time
}

type Response struct {
	Prices []PriceData
}

type Warning struct {
	Starts time.Time
	Ends   time.Time
}

type Args struct {
}

func main() {
	var args struct {
		Threshold int    `arg:"-t,--threshold,required" help:"Price threshold as ¢/kWh"`
		GotifyURL string `arg:"-g,--gotify,required" help:"Gotify URL"`
		GotifyKey string `arg:"-k,--key,required" help:"Gotify API key"`
		Lang      string `arg:"-l,--lang" help:"Notification language [fi, en]. Default: fi"`
	}
	arg.MustParse(&args)

	gotifyURL := fmt.Sprintf("%s/message?token=%s", strings.TrimSuffix(args.GotifyURL, "/"), url.QueryEscape(args.GotifyKey))
	apiUrl := "https://api.porssisahko.net/v1/latest-prices.json"

	priceClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Fatalf("creating the price request failed: %v", err)
	}

	res, err := priceClient.Do(req)
	if err != nil {
		log.Fatalf("making the price request failed: %v", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("reading the response body failed: %v", err)
	}

	response := Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("unmarshaling the response body failed: %v", err)
	}

	for i, j := 0, len(response.Prices)-1; i < j; i, j = i+1, j-1 {
		response.Prices[i], response.Prices[j] = response.Prices[j], response.Prices[i]
	}

	tomorrow := time.Now().Add(time.Hour * 24)
	tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.Local)
	warnings := make([]Warning, 0)

	var startTime time.Time
	var startSet bool = false
	for i, v := range response.Prices {
		if v.StartDate.Local().Before(tomorrow) {
			continue
		}
		if v.StartDate.After(tomorrow) && v.Price > float32(args.Threshold) && !startSet {
			startTime = time.Date(v.StartDate.Year(), v.StartDate.Month(), v.StartDate.Day(), v.StartDate.Hour(), v.StartDate.Minute(), v.StartDate.Second(), 0, v.StartDate.Location())
			startSet = true
		} else if (v.Price <= float32(args.Threshold) || i == len(response.Prices)-1) && startSet {
			warnings = append(warnings, Warning{Starts: startTime, Ends: v.StartDate})
			startSet = false
		}
	}

	for _, v := range warnings {
		message := fmt.Sprintf("Yli hintarajan (%d ¢/kWh) kello %s. Alle hintarajan kello %s.", args.Threshold, v.Starts.Local().Format("15:04"), v.Ends.Local().Format("15:04"))
		title := "Sähkön hinta"

		if strings.ToLower(args.Lang) == "en" {
			message = fmt.Sprintf("Above the price threshold (%d ¢/kWh) at %s. Below the price threshold at %s.", args.Threshold, v.Starts.Local().Format("15:04"), v.Ends.Local().Format("15:04"))
			title = "Electricity price"
		}

		resp, err := http.PostForm(gotifyURL, url.Values{"message": {message}, "title": {title}})
		if err != nil {
			log.Fatalf("failed to send the notication: %v", err)
		}

		if resp.StatusCode > http.StatusOK {
			log.Fatalf("failed to send the notication: server responded with status code %s", resp.Status)
		}

		fmt.Println(message)
	}
}
