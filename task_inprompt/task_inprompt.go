package task_inprompt

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type InpromptTask struct {
	Code     int64    `json:"code"`
	Msg      string   `json:"msg"`
	Input    []string `json:"input"`
	Question string   `json:"question"`
}

func TaskInprompt(apiKey string) {
	taskName := "inprompt"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))
	fmt.Print("\n\n")

	var result InpromptTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	name, err := extractNameFromQuestion(result.Question)
	if err != nil {
		panic(err)
	}

	relevantInputs := []string{}
	for _, input := range result.Input {
		if strings.Contains(input, name) {
			relevantInputs = append(relevantInputs, input)
		}
	}

	for _, sentence := range relevantInputs {
		answer, err := answerQuestion(result.Question, sentence)
		if err != nil {
			panic(err)
		}
		fmt.Println(answer)
		answerResult := ailabs.AnswerString(token, answer)
		fmt.Println(answerResult)
	}
}

func extractNameFromQuestion(question string) (string, error) {

	messages := []openapi.Message{
		{Role: "system", Content: "I return only a person's name from the given sentence."},
		{Role: "user", Content: question},
	}
	response := openapi.GetCompletion(messages, "gpt-4")
	if len(response.Choices) == 0 {
		return "", errors.New("could not find name")
	}

	return response.Choices[0].Message.Content, nil
}

func answerQuestion(question string, facts string) (string, error) {
	messages := []openapi.Message{
		{Role: "system", Content: facts},
		{Role: "user", Content: question},
	}
	response := openapi.GetCompletion(messages, "gpt-4")
	if len(response.Choices) == 0 {
		return "", errors.New("could not answer question")
	}

	return response.Choices[0].Message.Content, nil
}
