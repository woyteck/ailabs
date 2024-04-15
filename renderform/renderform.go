package renderform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Data struct {
	Text  string `json:"text.text"`
	Image string `json:"image.src"`
}

type RenderRequest struct {
	Template string `json:"template"`
	Data     Data   `json:"data"`
	FileName string `json:"fileName"`
	Version  string `json:"version"`
}

type RenderResponse struct {
	RequestId string `json:"requestId"`
	Href      string `json:"href"`
}

func Render(request RenderRequest) RenderResponse {
	url := "https://api.renderform.io/api/v2/render"

	postBody, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}
	req.Header.Add("x-api-key", os.Getenv("RENDERFORM_KEY"))
	req.Header.Add("Content-Type", "application/json")

	fmt.Println(string(postBody))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()
	if response.StatusCode >= 400 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Coult not read response")
		}
		fmt.Println(string(body))
	}

	var result RenderResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	return result
}
