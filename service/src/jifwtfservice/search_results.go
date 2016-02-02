package main

type SearchResults struct {
	Query         string         `json:"query"`
	SearchResults []SearchResult `json:"search_results"`
}
