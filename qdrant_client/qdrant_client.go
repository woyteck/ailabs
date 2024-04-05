package qdrant_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Qdrant struct {
}

type Point struct {
	Id      int            `json:"id"`
	Vector  []float64      `json:"vector"`
	Payload map[string]any `json:"payload"`
}

type UpsertPointsRequest struct {
	Points []Point `json:"points"`
}

type UpsertPointsResult struct {
	OperationId int    `json:"operation_id"`
	Stauts      string `json:"status"`
}

type UpsertPointsResponse struct {
	Result UpsertPointsResult `json:"result"`
	Status string             `json:"status"`
	Time   float64            `json:"time"`
}

type SearchRequest struct {
	Vector      []float64 `json:"vector"`
	Top         int       `json:"top"`
	WithPayload bool      `json:"with_payload"`
}

type SearchResult struct {
	Id      int            `json:"id"`
	Score   float64        `json:"score"`
	Payload map[string]any `json:"payload"`
	Version int            `json:"version"`
}

type SearchResponse struct {
	Result []SearchResult `json:"result"`
	Status string         `json:"status"`
	Time   float64        `json:"time"`
}

func NewClient() *Qdrant {
	return &Qdrant{}
}

func (qdrant *Qdrant) UpsertPoints(collectionName string, vector []float64, id int, payload map[string]any) UpsertPointsResponse {
	url := fmt.Sprintf("http://localhost:6333/collections/%v/points?wait=true", collectionName)

	request := UpsertPointsRequest{
		Points: []Point{
			{
				Id:      id,
				Vector:  vector,
				Payload: payload,
			},
		},
	}

	body, _ := json.Marshal(request)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()

	var result UpsertPointsResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	return result
}

func (qdrant *Qdrant) Search(collectionName string, vector []float64, resultsCount int) SearchResponse {
	url := fmt.Sprintf("http://localhost:6333/collections/%v/points/search", collectionName)

	request := SearchRequest{
		Vector:      vector,
		Top:         resultsCount,
		WithPayload: true,
	}

	body, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()

	var result SearchResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	return result
}
