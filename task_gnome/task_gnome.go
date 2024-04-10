package task_gnome

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type GnomeTask struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Hint string `json:"hint"`
	Url  string `json:"url"`
}

func TaskGnome(apiKey string) {
	taskName := "gnome"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result GnomeTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	answer, err := GetCompletion("If there is a gnome on this picture, tell me his hat's color in polish. Return only the color, nothing else. I you can't do it just return the word ERROR", result.Url)
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)

	answerResult := ailabs.AnswerString(token, answer)
	fmt.Println(answerResult)
}

func GetCompletion(prompt string, url string) (string, error) {
	messages := []openapi.ImageMessage{
		{Role: "user", Content: []openapi.Content{
			{Type: "text", Text: prompt},
			{Type: "image_url", ImageURL: openapi.ImageURL{
				URL: url,
			}},
		}},
	}

	response := openapi.GetImageCompletionShort(messages, "gpt-4-turbo")
	if len(response.Choices) == 0 {
		return "", errors.New("no choices")
	}

	return response.Choices[0].Message.Content, nil
}
