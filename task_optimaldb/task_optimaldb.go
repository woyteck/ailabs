package task_optimaldb

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type OptimaldbTask struct {
	Code     int64  `json:"code"`
	Msg      string `json:"msg"`
	Database string `json:"database"`
	Hint     string `json:"hint"`
}

func TaskOptimaldb(apiKey string) {
	taskName := "optimaldb"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result OptimaldbTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	compressed := compressDatabase(result.Database)

	fmt.Println(compressed)

	answerResult := ailabs.AnswerString(token, compressed)
	fmt.Println(answerResult)
}

func compressDatabase(url string) string {
	bytes, _ := downloadFile(url)

	var inputDatabase = map[string][]string{}
	err := json.Unmarshal(bytes, &inputDatabase)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	result := ""
	for personName, facts := range inputDatabase {
		compressed := compressPersonFacts(facts)

		result += personName
		result += ":\n"
		result += compressed
		result += "\n"
	}

	return result
}

func compressPersonFacts(facts []string) string {
	prompt := ""
	for _, fact := range facts {
		prompt += fact
		prompt += "\n"
	}

	context := "Skracam tekst zachowując wszystkie informacje.\n"
	context += "Staram się skrócić tekst o 1/3\n"
	context += "Nie powtarzam imienia osoby, użytkownik już je zna"

	messages := []openapi.Message{
		{
			Role:    "system",
			Content: context,
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}
	response := openapi.GetCompletionShort(messages, "gpt-4-turbo")
	text := response.Choices[0].Message.Content

	return text
}

func downloadFile(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
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
