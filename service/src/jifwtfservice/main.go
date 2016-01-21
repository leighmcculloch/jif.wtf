package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/googollee/go-socket.io"
)

func onConnection(so socketio.Socket) {
	so.On("search", func(engine, query string) {
		var results []SearchResult
		switch engine {
		case "imgur":
			results = GetImgurResults(query)
		case "giphy":
			results = GetGiphyResults(query)
		}
		sort.Sort(SearchResultsByRating(results))

		so.Emit("results", query, results)
	})
}

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.SetCookie("JSESSIONID")

	server.On("connection", onConnection)

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
