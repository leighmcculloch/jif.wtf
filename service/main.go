package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/search", search)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

const tenorBaseURL = "https://api.tenor.com/v1/search"

var tenorKey = os.Getenv("TENOR_KEY")

func search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	tenorParams := url.Values{}
	tenorParams.Set("key", tenorKey)
	tenorParams.Set("safesearch", "mild")
	tenorParams.Set("limit", "25")
	tenorParams.Set("media_filter", "minimal")
	tenorParams.Set("q", q)
	tenorURL := tenorBaseURL + "?" + tenorParams.Encode()

	req, err := http.NewRequest("GET", tenorURL, nil)
	if err != nil {
		panic(err)
	}
	reqContext, reqCancel := context.WithTimeout(r.Context(), time.Minute)
	defer reqCancel()
	req = req.WithContext(reqContext)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var tenorRes struct {
		Results []struct {
			Media []struct {
				GIF struct {
					URL string `json:"url"`
				} `json:"gif"`
				MP4 struct {
					URL        string `json:"url"`
					Dimensions []int  `json:"dims"`
				} `json:"mp4"`
			} `json:"media"`
		} `json:"results"`
	}
	err = json.Unmarshal(body, &tenorRes)
	if err != nil {
		panic(err)
	}

	type ResponseDataImagesOriginal struct {
		URL    string `json:"url"`
		MP4    string `json:"mp4"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}
	type ResponseDataImages struct {
		Original ResponseDataImagesOriginal `json:"original"`
	}
	type ResponseData struct {
		Images ResponseDataImages `json:"images"`
	}
	type Response struct {
		Data []ResponseData `json:"data"`
	}

	res := Response{
		Data: []ResponseData{},
	}
	for _, r := range tenorRes.Results {
		res.Data = append(res.Data, ResponseData{
			Images: ResponseDataImages{
				Original: ResponseDataImagesOriginal{
					URL:    r.Media[0].GIF.URL,
					MP4:    r.Media[0].MP4.URL,
					Width:  r.Media[0].MP4.Dimensions[0],
					Height: r.Media[0].MP4.Dimensions[1],
				},
			},
		})
	}

	resBytes, err := json.MarshalIndent(res, "", "  ")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	_, err = w.Write(resBytes)
	if err != nil {
		panic(err)
	}
}
