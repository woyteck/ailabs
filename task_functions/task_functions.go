package task_functions

import (
	"encoding/json"
	"fmt"
	"log"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type FunctionsTask struct {
	Code  int64  `json:"code"`
	Msg   string `json:"msg"`
	Hint1 string `json:"hint1"`
}

type AddUserProperties struct {
	Name    openapi.Param `json:"name"`
	Surname openapi.Param `json:"surname"`
	Year    openapi.Param `json:"year"`
}

type Parameters struct {
	Type       string            `json:"type"`
	Properties AddUserProperties `json:"properties"`
}

func TaskFunctions(apiKey string) {
	taskName := "functions"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result FunctionsTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	request := openapi.CompletionRequest{
		Model: "gpt-4",
		Messages: []openapi.Message{
			{Role: "system", Content: ""},
			{Role: "user", Content: "Hello"},
		},
		Tools: []openapi.Tool{
			{
				Type: "function",
				Function: openapi.Function{
					Name:        "addUser",
					Description: "Adds user",
					Parameters: Parameters{
						Type: "object",
						Properties: AddUserProperties{
							Name: openapi.Param{
								Type:        "string",
								Description: "User's name",
							},
							Surname: openapi.Param{
								Type:        "string",
								Description: "User's surname",
							},
							Year: openapi.Param{
								Type:        "number",
								Description: "User's year of birth",
							},
						},
					},
				},
			},
		},
	}

	postBody, _ := json.Marshal(request.Tools[0].Function)
	fmt.Println(string(postBody))

	answerResult := ailabs.AnswerAny(token, request.Tools[0].Function)
	fmt.Println(answerResult)
}
