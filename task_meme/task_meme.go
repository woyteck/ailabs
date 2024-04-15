package task_meme

import (
	"encoding/json"
	"fmt"
	"log"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/renderform"
)

type MemeTask struct {
	Code    int64  `json:"code"`
	Msg     string `json:"msg"`
	Service string `json:"service"`
	Image   string `json:"image"`
	Text    string `json:"text"`
	Hint    string `json:"hint"`
}

func TaskMeme(apiKey string) {
	taskName := "meme"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result MemeTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	request := renderform.RenderRequest{
		Template: "faulty-yetis-attack-soon-1236",
		Data: renderform.Data{
			Text:  result.Text,
			Image: result.Image,
		},
		FileName: "meme",
		Version:  "v3",
	}
	response := renderform.Render(request)
	fmt.Println(response.Href)

	answerResult := ailabs.AnswerString(token, response.Href)
	fmt.Println(answerResult)
}
