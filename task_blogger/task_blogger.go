package task_blogger

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type BlogTask struct {
	Code int64    `json:"code"`
	Msg  string   `json:"msg"`
	Blog []string `json:"blog"`
}

func TaskBlogger(apiKey string) {
	taskName := "blogger"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result BlogTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	context := "I'm a cooking blogger."
	prompt := "Write a blog article about making pizza Margherita, write 4 chapters titled:\n"
	for _, result := range result.Blog {
		prompt += result + "\n"
	}
	prompt += "\nSeparate chapters with this line: ####"

	messages := []openapi.Message{
		{Role: "system", Content: context},
		{Role: "user", Content: prompt},
	}
	response := openapi.GetCompletionShort(messages, "gpt-4")

	text := response.Choices[0].Message.Content

	answerResult := ailabs.AnswerAny(token, strings.Split(text, "####"))
	fmt.Println(answerResult)
}
