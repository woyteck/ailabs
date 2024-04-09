package task_tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type ToolsTask struct {
	Code               int64  `json:"code"`
	Msg                string `json:"msg"`
	Hint               string `json:"hint"`
	ExampleForTodo     string `json:"example for ToDo"`
	ExampleForCalendar string `json:"example for Calendar"`
	Question           string `json:"question"`
}

type Message struct {
	Tool string `json:"tool"`
	Desc string `json:"desc"`
	Date string `json:"date,omitempty"`
}

func TaskTools(apiKey string) {
	taskName := "tools"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result ToolsTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	answer, _ := decideToolToUse(result.Question, result.ExampleForTodo, result.ExampleForCalendar)
	fmt.Println(answer)

	var message Message
	err = json.Unmarshal([]byte(answer), &message)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	answerResult := ailabs.AnswerAny(token, message)
	fmt.Println(answerResult)
}

func decideToolToUse(query string, todoExample string, calendarExample string) (string, error) {
	now := time.Now()

	context := "I choose which tool to use based on given query\n"
	context += "Tools available: ToDo, calendar\n"
	context += "If I can determine the date i should return the calendar tool. Today is "
	context += now.Format(time.DateOnly)
	context += "\n"
	context += "I return responses only in json format:\n"
	context += "example for ToDo:\n"
	context += todoExample
	context += "\nexample for calendar:\n"
	context += calendarExample

	prompt := query

	messages := []openapi.Message{
		{Role: "system", Content: context},
		{Role: "user", Content: prompt},
	}
	response := openapi.GetCompletionShort(messages, "gpt-4")

	if len(response.Choices) == 0 {
		return "", errors.New("no choices returned")
	}

	text := response.Choices[0].Message.Content

	return text, nil
}
