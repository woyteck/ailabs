package ailabs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type GetTokenResponse struct {
	Code  int64  `json:"code"`
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

func GetToken(taskName string, apiKey string) string {
	url := fmt.Sprintf("https://tasks.aidevs.pl/token/%v", taskName)
	postBody, _ := json.Marshal(map[string]string{
		"apikey": apiKey,
	})
	body := bytes.NewBuffer(postBody)
	response, err := http.Post(url, "application/json", body)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()

	var result GetTokenResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	return result.Token
}

func GetTask(token string) []byte {
	url := fmt.Sprintf("https://tasks.aidevs.pl/task/%v", token)
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Coult not read response")
	}

	return body
}

func AnswerString(token string, answer string) string {
	url := fmt.Sprintf("https://tasks.aidevs.pl/answer/%v", token)
	postBody, _ := json.Marshal(map[string]string{
		"answer": answer,
	})
	payload := bytes.NewBuffer(postBody)
	response, err := http.Post(url, "application/json", payload)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Coult not read response")
	}

	return string(body)
}

func AnswerAny[T any](token string, answer T) string {
	url := fmt.Sprintf("https://tasks.aidevs.pl/answer/%v", token)
	postBody, _ := json.Marshal(map[string]T{
		"answer": answer,
	})
	payload := bytes.NewBuffer(postBody)
	response, err := http.Post(url, "application/json", payload)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Coult not read response")
	}

	return string(body)
}
