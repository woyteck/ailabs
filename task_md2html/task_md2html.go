package task_md2html

import (
	"encoding/json"
	"fmt"
	"log"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type Md2htmlTask struct {
	Code  int64  `json:"code"`
	Msg   string `json:"msg"`
	Hint  string `json:"hint"`
	Input string `json:"input"`
}

func TaskMd2html(apiKey string) {
	taskName := "md2html"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result Md2htmlTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	fmt.Println(result.Input)
	fmt.Println("")
	fmt.Println("")

	messages := []openapi.Message{
		{Role: "system", Content: "I convert given text to html"},
		{Role: "user", Content: result.Input},
	}
	completions := openapi.GetCompletionShort(messages, "ft:gpt-3.5-turbo-0125:personal::secred-well-finetuned-model")
	if len(completions.Choices) == 0 {
		panic("no completions returned")
	}

	fmt.Println(completions.Choices[0].Message.Content)

	answerResult := ailabs.AnswerString(token, completions.Choices[0].Message.Content)
	fmt.Println(answerResult)
}
