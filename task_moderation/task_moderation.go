package task_moderation

import (
	"encoding/json"
	"fmt"
	"log"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type ModerationTask struct {
	Code  int64    `json:"code"`
	Msg   string   `json:"msg"`
	Input []string `json:"input"`
}

func TaskModeration(apiKey string) {
	taskName := "moderation"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result ModerationTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	results := []int{}
	for _, result := range result.Input {
		isFlagged, _ := openapi.GetModeration(result)
		if isFlagged {
			results = append(results, 1)
		} else {
			results = append(results, 0)
		}
	}
	fmt.Println(results)

	answerResult := ailabs.AnswerAny(token, results)
	fmt.Println(answerResult)
}
