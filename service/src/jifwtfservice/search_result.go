package main

type SearchResult struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Mp4    string `json:"mp4"`
	Webp   string `json:"webp"`
	Gif    string `json:"gif"`
	Rating int    `json:"rating"`
}

type SearchResultsByRating []SearchResult

func (s SearchResultsByRating) Len() int {
	return len(s)
}

func (s SearchResultsByRating) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SearchResultsByRating) Less(i, j int) bool {
	return s[i].Rating > s[j].Rating
}
