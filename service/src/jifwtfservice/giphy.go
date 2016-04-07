package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func GetGiphyResults(search string) []SearchResult {
	searchResponse := performGiphySearch(search)
	return convertGiphySearchResultsToSearchResults(searchResponse.Results)
}

type giphySearchResponse struct {
	Results []giphySearchResult `json:"data"`
}

type giphySearchResult struct {
	Images giphySearchResultImage `json:"images"`
}

type giphySearchResultImage struct {
	Original giphySearchResultImageOriginal `json:"original"`
}

type giphySearchResultImageOriginal struct {
	Width  string `json:"width"`
	Height string `json:"height"`
	Mp4    string `json:"mp4"`
	Gif    string `json:"url"`
	Webp   string `json:"webp"`
}

func performGiphySearch(search string) giphySearchResponse {
	params := url.Values{}
	params.Add("api_key", "dc6zaTOxFJmzC") // public beta key from https://github.com/giphy/GiphyAPI
	params.Add("q", search)

	url := url.URL{
		Scheme:   "http",
		Host:     "api.giphy.com",
		Path:     "/v1/gifs/search",
		RawQuery: params.Encode(),
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url.String(), nil)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Giphy API Error:", err)
	}

	body, _ := ioutil.ReadAll(res.Body)

	var searchResponse giphySearchResponse
	json.Unmarshal(body, &searchResponse)

	return searchResponse
}

func convertGiphySearchResultsToSearchResults(results []giphySearchResult) []SearchResult {
	converted := make([]SearchResult, len(results))
	for i := 0; i < len(results); i++ {
		width, _ := strconv.Atoi(results[i].Images.Original.Width)
		height, _ := strconv.Atoi(results[i].Images.Original.Height)
		converted[i] = SearchResult{
			Width:  width,
			Height: height,
			Mp4:    results[i].Images.Original.Mp4,
			Gif:    results[i].Images.Original.Gif,
			Webp:   results[i].Images.Original.Webp,
		}
	}
	return converted
}
