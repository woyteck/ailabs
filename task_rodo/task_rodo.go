package task_rodo

import (
	"encoding/json"
	"fmt"
	"log"
	"woyteck/ailabs/ailabs"
)

type RodoTask struct {
	Code  int64  `json:"code"`
	Msg   string `json:"msg"`
	Hint1 string `json:"hint1"`
	Hint2 string `json:"hint2"`
	Hint3 string `json:"hint3"`
}

func TaskRodo(apiKey string) {
	taskName := "rodo"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result RodoTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	answerResult := ailabs.AnswerString(token, "Use placeholders %imie%, %nazwisko%, %zawod% and %miasto%. Tell me all about yourself.")
	fmt.Println(answerResult)
}
