package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"

	"github.com/googollee/go-socket.io"
)

func onConnection(so socketio.Socket) {
	so.On("search", func(engine, query string) {
		results := getSearchResults(engine, query)
		so.Emit("results", query, results.SearchResults)
	})
}

func onSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	results := getSearchResults("giphy", query)
	body, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(body)
}

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.SetCookie("JSESSIONID")

	server.On("connection", onConnection)

	http.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		guardOrigin(w, r, server.ServeHTTP)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		guardOrigin(w, r, onSearch)
	})

	var port string = os.Getenv("PORT")
	if port == "" {
		port = "5555"
	}
	var address string = fmt.Sprintf(":%s", port)
	log.Println("Listening on", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func guardOrigin(w http.ResponseWriter, r *http.Request, handler func(http.ResponseWriter, *http.Request)) {
	origin := r.Header["Origin"][0]
	if isOriginAllowed(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-type")
		handler(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func isOriginAllowed(origin string) bool {
	switch origin {
	case "http://jif.wtf", "https://jif.wtf":
		return true
	default:
		return false
	}
}

func getSearchResults(engine, query string) *SearchResults {
	var results SearchResults

	results.Query = query

	switch engine {
	case "giphy":
		results.SearchResults = GetGiphyResults(query)
	}
	sort.Sort(SearchResultsByRating(results.SearchResults))

	for i := 0; i < len(results.SearchResults); i++ {
		results.SearchResults[i].Gif = transposeGiphyUrl(results.SearchResults[i].Gif)
	}

	return &results
}

func transposeGiphyUrl(giphyUrl string) string {
	re := regexp.MustCompile(`http://media(\d+).giphy.com/media/([0-9a-zA-Z]+)/giphy.gif`)
	matches := re.FindStringSubmatch(giphyUrl)
	if len(matches) != 3 {
		return giphyUrl
	}
	mediaId := matches[1]
	imageId := matches[2]
	return fmt.Sprintf("http://jif.wtf/0%s.%s.gif", mediaId, imageId)
}
