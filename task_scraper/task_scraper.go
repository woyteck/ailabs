package task_scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type ScraperTask struct {
	Code     int64  `json:"code"`
	Msg      string `json:"msg"`
	Input    string `json:"input"`
	Question string `json:"question"`
}

func TaskScraper(apiKey string) {
	taskName := "scraper"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result ScraperTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	article, err := getArticle(result.Input)
	if err != nil {
		panic(err)
	}

	context := article
	prompt := result.Question

	messages := []openapi.Message{
		{Role: "system", Content: context},
		{Role: "user", Content: prompt},
	}
	response := openapi.GetCompletionShort(messages, "gpt-4")

	answer := response.Choices[0].Message.Content

	answerResult := ailabs.AnswerString(token, answer)
	fmt.Println(answerResult)
}

func getArticle(url string) (string, error) {
	var articleString string
	maxTriesCount := 3
	triesCount := 0
	for {
		if triesCount > maxTriesCount {
			return "", errors.New("maximum try count exceeded")
		}

		article, err := downloadFile(url)
		if err == nil {
			articleString = string(article)
			break
		} else {
			time.Sleep(time.Second)
		}
		triesCount += 1
	}

	return articleString, nil
}

func downloadFile(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
