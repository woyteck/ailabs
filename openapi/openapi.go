package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Function struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Param struct {
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	N        int       `json:"n,omitempty"`
	Stream   bool      `json:"stream,omitempty"`
	User     string    `json:"user,omitempty"`
	Tools    []Tool    `json:"tools"`
}

type CompletionResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type ModerationCategories struct {
	Sexual                bool `json:"sexual"`
	Hate                  bool `json:"hate"`
	Harassment            bool `json:"harassment"`
	SelfHarm              bool `json:"self-harm"`
	SexualMinors          bool `json:"sexual/minors"`
	HateThreatening       bool `json:"hate/threatening"`
	ViolenceGraphic       bool `json:"violence/graphic"`
	SelfHarmIntent        bool `json:"self-harm/intent"`
	SelfHarmInstructions  bool `json:"self-harm/instructions"`
	HarassmentThreatening bool `json:"harassment/threatening"`
	Violence              bool `json:"violence"`
}

type ModerationCategoryScores struct {
	Sexual                float64 `json:"sexual"`
	Hate                  float64 `json:"hate"`
	Harassment            float64 `json:"harassment"`
	SelfHarm              float64 `json:"self-harm"`
	SexualMinors          float64 `json:"sexual/minors"`
	HateThreatening       float64 `json:"hate/threatening"`
	ViolenceGraphic       float64 `json:"violence/graphic"`
	SelfHarmIntent        float64 `json:"self-harm/intent"`
	SelfHarmInstructions  float64 `json:"self-harm/instructions"`
	HarassmentThreatening float64 `json:"harassment/threatening"`
	Violence              float64 `json:"violence"`
}

type ModerationResult struct {
	Flagged        bool                     `json:"flagged"`
	Categories     ModerationCategories     `json:"categories"`
	CategoryScores ModerationCategoryScores `json:"category_scores"`
}

type ModerationRequest struct {
	Input string `json:"input"`
}

type ModerationResponse struct {
	Id      string             `json:"id"`
	Model   string             `json:"model"`
	Results []ModerationResult `json:"results"`
}

type EmbeddingRequest struct {
	Input          string `json:"input"`
	Model          string `json:"model"`
	EncodingFormat string `json:"encoding_format"`
}

type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  Usage           `json:"usage"`
}

type TranscriptionResponse struct {
	Text string `json:"text"`
}

func GetCompletion(request CompletionRequest) CompletionResponse {
	url := "https://api.openai.com/v1/chat/completions"

	postBody, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("OPENAI_API_KEY")))
	req.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()
	if response.StatusCode >= 400 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Coult not read response")
		}
		fmt.Println(string(body))
	}

	var result CompletionResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	return result
}

func GetCompletionShort(messages []Message, model string) CompletionResponse {
	request := CompletionRequest{
		Model:    model,
		Messages: messages,
	}

	return GetCompletion(request)
}

func GetModeration(input string) (bool, ModerationResponse) {
	url := "https://api.openai.com/v1/moderations"

	request := ModerationRequest{
		Input: input,
	}

	postBody, _ := json.Marshal(request)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("OPENAI_API_KEY")))
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()

	var result ModerationResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	isFlagged := false
	for _, result := range result.Results {
		if result.Flagged {
			isFlagged = true
		}
	}

	return isFlagged, result
}

func GetEmbedding(input string, model string) []float64 {
	url := "https://api.openai.com/v1/embeddings"

	request := EmbeddingRequest{
		Input:          input,
		Model:          model,
		EncodingFormat: "float",
	}

	postBody, _ := json.Marshal(request)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("OPENAI_API_KEY")))
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()

	var result EmbeddingResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	if len(result.Data) == 0 {
		log.Fatal("No embeddings returned")
	}

	return result.Data[0].Embedding
}

func GetTranscription(file []byte, model string) string {
	url := "https://api.openai.com/v1/audio/transcriptions"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	filePart, _ := writer.CreateFormFile("file", "file.mp3")
	filePart.Write(file)
	writer.WriteField("model", model)
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("OPENAI_API_KEY")))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	fmt.Println(response.StatusCode)

	defer response.Body.Close()

	var result TranscriptionResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	return result.Text
}
