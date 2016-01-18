package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"

	"github.com/googollee/go-socket.io"
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.SetCookie("JSESSIONID")

	server.On("connection", func(so socketio.Socket) {
		so.On("search", func(query string) {
			log.Println("search:", query)

			body := getResults(query)

			so.Emit("results", query, body)
		})
	})

	http.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://jif.wtf")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-type")
		server.ServeHTTP(w, r)
	})

	var host string
	var port string = os.Getenv("PORT")
	if port == "" {
		host = "127.0.0.1"
		port = "5555"
	}
	var address string = fmt.Sprintf("%s:%s", host, port)
	log.Println("Listening on", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

type SearchResult struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Mp4       string `json:"mp4"`
	Gifv      string `json:"gifv"`
	Gif       string `json:"link"`
	IsLooping bool   `json:"looping"`
	Points    int    `json:"points"`
}

type SearchResponse struct {
	Success bool           `json:"success"`
	Results []SearchResult `json:"data"`
}

func getResults(search string) []SearchResult {
	params := url.Values{}
	params.Add("q_type", "anigif")
	params.Add("q_all", search)

	url := url.URL{
		Scheme:   "https",
		Host:     "api.imgur.com",
		Path:     "/3/gallery/search",
		RawQuery: params.Encode(),
	}
	log.Println("url:", url.String())

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Set("Authorization", "Client-ID a27edb3db737edf")
	res, err := client.Do(req)
	if err != nil {
		log.Println("imgur api error:", err)
	}

	body, _ := ioutil.ReadAll(res.Body)

	var searchResponse SearchResponse
	json.Unmarshal(body, &searchResponse)

	results := searchResponse.Results
	sort.Sort(ByPoints(results))
	return results
}

type ByPoints []SearchResult

func (s ByPoints) Len() int {
	return len(s)
}
func (s ByPoints) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByPoints) Less(i, j int) bool {
	return s[i].Points < s[j].Points
}
