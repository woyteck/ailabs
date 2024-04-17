package task_google

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/serp"

	"github.com/gin-gonic/gin"
)

type GoogleTask struct {
	Code     int64  `json:"code"`
	Msg      string `json:"msg"`
	Database string `json:"database"`
	Hint     string `json:"hint"`
}

type Answer struct {
	Reply string `json:"reply"`
}

type Question struct {
	Question string `json:"question"`
}

func TaskGoogle(apiKey string) {
	taskName := "google"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result GoogleTask
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

	results := serp.Search(question.Question)
	if len(results.OrganicResults) == 0 {
		panic("no results")
	}

	response := Answer{
		Reply: results.OrganicResults[0].Link,
	}

	context.JSON(http.StatusOK, response)
}
