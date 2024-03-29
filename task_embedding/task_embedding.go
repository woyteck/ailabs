package task_embedding

import (
	"encoding/json"
	"fmt"
	"log"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type EmbeddingTask struct {
	Code  int64  `json:"code"`
	Msg   string `json:"msg"`
	Hint1 string `json:"hint1"`
	Hint2 string `json:"hint2"`
	Hint3 string `json:"hint3"`
}

func TaskEmbedding(apiKey string) {
	taskName := "embedding"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result EmbeddingTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	embeddings := openapi.GetEmbedding("Hawaiian pizza", "text-embedding-ada-002")

	answerResult := ailabs.AnswerAny(token, embeddings[0].Embedding)
	fmt.Println(answerResult)
}
