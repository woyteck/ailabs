package task_helloapi

import (
	"encoding/json"
	"fmt"
	"log"
	"woyteck/ailabs/ailabs"
)

type HelloApiTask struct {
	Code   int64  `json:"code"`
	Msg    string `json:"msg"`
	Cookie string `json:"cookie"`
}

func TaskHelloApi(apiKey string) {
	taskName := "helloapi"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result HelloApiTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	answerResult := ailabs.AnswerString(token, result.Cookie)
	fmt.Println(answerResult)
}
