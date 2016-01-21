package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func GetImgurResults(search string) []SearchResult {
	searchResponse := performImgurSearch(search)
	return convertImgurSearchResultsToSearchResults(searchResponse.Results)
}

type imgurSearchResponse struct {
	Success bool                `json:"success"`
	Results []imgurSearchResult `json:"data"`
}

type imgurSearchResult struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Mp4       string `json:"mp4"`
	Gifv      string `json:"gifv"`
	Gif       string `json:"link"`
	IsLooping bool   `json:"looping"`
	Points    int    `json:"points"`
}

func performImgurSearch(search string) imgurSearchResponse {
	params := url.Values{}
	params.Add("q_type", "anigif")
	params.Add("q_all", search)

	url := url.URL{
		Scheme:   "https",
		Host:     "api.imgur.com",
		Path:     "/3/gallery/search",
		RawQuery: params.Encode(),
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Set("Authorization", "Client-ID a27edb3db737edf")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Imgur API Error:", err)
	}

	body, _ := ioutil.ReadAll(res.Body)

	var searchResponse imgurSearchResponse
	json.Unmarshal(body, &searchResponse)

	return searchResponse
}

func convertImgurSearchResultsToSearchResults(results []imgurSearchResult) []SearchResult {
	converted := make([]SearchResult, len(results))
	for i := 0; i < len(results); i++ {
		converted[i] = SearchResult{
			Width:  results[i].Width,
			Height: results[i].Height,
			Mp4:    results[i].Mp4,
			Gif:    results[i].Gif,
			Rating: results[i].Points,
		}
	}
	return converted
}
