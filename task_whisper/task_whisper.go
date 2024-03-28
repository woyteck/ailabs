package task_whisper

import (
	"fmt"
	"io"
	"net/http"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

func TaskWhisper(apiKey string) {
	taskName := "whisper"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	soundUrl := "https://tasks.aidevs.pl/data/mateusz.mp3"
	file, err := downloadFile(soundUrl)
	if err != nil {
		panic(err)
	}

	text := openapi.GetTranscription(file, "whisper-1")
	fmt.Println(text)

	answerResult := ailabs.AnswerString(token, text)
	fmt.Println(answerResult)
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
