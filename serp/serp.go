package serp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type OrganicResult struct {
	Position int    `json:"position"`
	Title    string `json:"title"`
	Link     string `json:"link"`
	Snippet  string `json:"snipper"`
}

type Results struct {
	OrganicResults []OrganicResult `json:"organic_results"`
}

func Search(query string) Results {
	key := os.Getenv("SERP_KEY")

	query = strings.ReplaceAll(query, " ", "+")

	params := url.Values{}
	params.Add("api_key", key)
	params.Add("q", query)

	url := "https://serpapi.com/search?" + params.Encode()

	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Coult not read response")
	}

	var result Results
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}

	return result
}
