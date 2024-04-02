package task_liar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type LiarTask struct {
	Code  int64  `json:"code"`
	Msg   string `json:"msg"`
	Hint1 string `json:"hint1"`
	Hint2 string `json:"hint2"`
	Hint3 string `json:"hint3"`
}

type Answer struct {
	Code   int64  `json:"code"`
	Msg    string `json:"msg"`
	Answer string `json:"answer"`
}

func TaskLiar(apiKey string) {
	taskName := "liar"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))
	fmt.Print("\n\n")

	var result LiarTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	question1 := "What is capital of Poland?"
	answer1 := PostTask(token, question1)

	fmt.Println(question1)
	fmt.Println(answer1)
	fmt.Printf("Hint1: %v\n", answer1)
	IsLying := IsLying(question1, answer1)
	fmt.Println(IsLying)

	var answer string
	if IsLying {
		answer = "NO"
	} else {
		answer = "YES"
	}
	ailabs.AnswerString(token, answer)
}

func IsLying(question string, answer string) bool {
	context := "I'm a lying detector."
	prompt := "Is this answer relevant to the question?\n"
	prompt += fmt.Sprintf("Question: %v\n", question)
	prompt += fmt.Sprintf("Answer: %v\n", answer)
	prompt += "Give me a simple yes/no answer"

	messages := []openapi.Message{
		{Role: "system", Content: context},
		{Role: "user", Content: prompt},
	}
	response := openapi.GetCompletionShort(messages, "gpt-4")

	text := response.Choices[0].Message.Content

	return strings.ToLower(text) == "no"
}

func PostTask(token string, question string) string {
	apiUrl := fmt.Sprintf("https://tasks.aidevs.pl/task/%v", token)

	var multipartBody bytes.Buffer
	writer := multipart.NewWriter(&multipartBody)
	writer.WriteField("question", question)
	writer.Close()

	postData := url.Values{}
	postData.Set("question", question)
	response, err := http.Post(apiUrl, writer.FormDataContentType(), &multipartBody)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()

	var result Answer
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	return result.Answer
}
