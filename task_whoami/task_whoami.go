package task_whoami

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type WhoamiTask struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Hint string `json:"hint"`
}

func TaskWhoami(apiKey string) {
	taskName := "whoami"

	var hints []string = []string{}

	for {
		token := ailabs.GetToken(taskName, apiKey)
		taskJsonString := ailabs.GetTask(token)
		fmt.Println(string(taskJsonString))

		var result WhoamiTask
		err := json.Unmarshal(taskJsonString, &result)
		if err != nil {
			log.Fatal("Can not unmarshall JSON")
		}

		hints := append(hints, result.Hint)
		prompt := strings.Join(hints, "\n")

		messages := []openapi.Message{
			{Role: "system", Content: "I'm guessing who the person is based on hints I receive. I return only his firstname and lastname, nothing else. If I don't know tho the person is I respond with exactly 'need another hint' in english"},
			{Role: "user", Content: prompt},
		}
		response := openapi.GetCompletionShort(messages, "gpt-4")
		text := response.Choices[0].Message.Content
		fmt.Println(text)

		if text != "need another hint" {
			answerResult := ailabs.AnswerString(token, text)
			fmt.Println(answerResult)
			break
		}
	}

}
