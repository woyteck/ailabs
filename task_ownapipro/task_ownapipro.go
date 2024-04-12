package task_ownapipro

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"

	"github.com/gin-gonic/gin"
)

type OwnapiproTask struct {
	Code  int64  `json:"code"`
	Msg   string `json:"msg"`
	Hint1 string `json:"hint1"`
	Hint2 string `json:"hint2"`
	Hint3 string `json:"hint3"`
}

type Answer struct {
	Reply string `json:"reply"`
}

type Question struct {
	Question string `json:"question"`
}

var messages = []openapi.Message{
	{Role: "system", Content: "I answer general questions. I keep my answers consice."},
}

func TaskOwnapipro(apiKey string) {
	taskName := "ownapipro"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result OwnapiproTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	messages1 := make(chan string)
	messages2 := make(chan string)

	go func() {
		server := gin.Default()
		server.SetTrustedProxies(nil)
		RegisterRoutes(server)
		server.Run(":8080")
	}()

	go func() {
		time.Sleep(time.Second * 4)
		answerResult := ailabs.AnswerString(token, "https://secret.location.com/answer")
		fmt.Println(answerResult)
	}()

	fmt.Println(<-messages1)
	fmt.Println(<-messages2)
}

func RegisterRoutes(server *gin.Engine) {
	server.POST("/answer", getAnswer)
}

func getAnswer(context *gin.Context) {
	var question Question
	err := context.ShouldBindJSON(&question)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	fmt.Println(question.Question)
	messages = append(messages, openapi.Message{Role: "user", Content: question.Question})

	completion := openapi.GetCompletionShort(messages, "gpt-4")
	if len(completion.Choices) == 0 {
		context.JSON(http.StatusInternalServerError, "")
	}

	response := Answer{
		Reply: completion.Choices[0].Message.Content,
	}

	fmt.Println(response.Reply)
	messages = append(messages, openapi.Message{Role: "assistant", Content: response.Reply})

	context.JSON(http.StatusOK, response)
}
