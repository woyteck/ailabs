package task_people

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
	"woyteck/ailabs/qdrant_client"
)

type PeopleTask struct {
	Code     int64  `json:"code"`
	Msg      string `json:"msg"`
	Data     string `json:"data"`
	Question string `json:"question"`
	Hint1    string `json:"hint1"`
	Hint2    string `json:"hint2"`
}

type Person struct {
	FirstName                  string `json:"imie"`
	LastName                   string `json:"nazwisko"`
	Age                        int    `json:"wiek"`
	About                      string `json:"o_mnie"`
	FavouriteCptBombaCharacter string `json:"ulubiona_postac_z_kapitana_bomby"`
	FavouriteSeries            string `json:"ulubiony_serial"`
	FavouriteMovie             string `json:"ulubiony_film"`
	FavouriteColor             string `json:"ulubiony_kolor"`
}

func TaskPeople(apiKey string) {
	taskName := "people"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result PeopleTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	qdrant := qdrant_client.NewClient()
	// indexer.CreateCollection("people")
	// people := fetchPeople()
	// indexPeople(qdrant, people)
	searchVector := openapi.GetEmbedding(result.Question, "text-embedding-ada-002")
	searchResults := qdrant.Search("people", searchVector, 1)
	if len(searchResults.Result) > 0 {
		fmt.Println(searchResults.Result[0].Payload)
		payload := searchResults.Result[0].Payload

		summarizedPerson := summarizePerson(
			payload["first_name"].(string),
			payload["last_name"].(string),
			int(payload["age"].(float64)),
			payload["favourite_cpt_bomba_character"].(string),
			payload["favourite_series"].(string),
			payload["favourite_movie"].(string),
			payload["favourite_color"].(string),
			payload["about"].(string),
		)
		messages := []openapi.Message{
			{Role: "system", Content: summarizedPerson},
			{Role: "user", Content: result.Question},
		}
		completion := openapi.GetCompletionShort(messages, "gpt-3.5-turbo-0125")
		fmt.Println(completion.Choices[0].Message.Content)

		answerResult := ailabs.AnswerString(token, completion.Choices[0].Message.Content)
		fmt.Println(answerResult)
	}
}

func summarizePerson(firstName string, lastName string, age int, favouriteCptBombaCharacter string, favouriteSeries string, favouriteMovie string, favouriteColor string, about string) string {
	return fmt.Sprintf(
		"Imię: %v, Nazwisko: %v, wiek: %v, ulubiona postać z kapitana bomby: %v, ulubiony serial: %v, ulubiony film: %v, ulubiony kolor: %v, opis: %v",
		firstName,
		lastName,
		age,
		favouriteCptBombaCharacter,
		favouriteSeries,
		favouriteMovie,
		favouriteColor,
		about,
	)
}

func indexPeople(qdrant qdrant_client.Qdrant, people []Person) {
	for index, person := range people {
		textToIndex := summarizePerson(
			person.FirstName,
			person.LastName,
			person.Age,
			person.FavouriteCptBombaCharacter,
			person.FavouriteSeries,
			person.FavouriteMovie,
			person.FavouriteColor,
			person.About,
		)
		vector := openapi.GetEmbedding(textToIndex, "text-embedding-ada-002")

		payload := map[string]any{}
		payload["first_name"] = person.FirstName
		payload["last_name"] = person.LastName
		payload["age"] = person.Age
		payload["about"] = person.About
		payload["favourite_cpt_bomba_character"] = person.FavouriteCptBombaCharacter
		payload["favourite_series"] = person.FavouriteSeries
		payload["favourite_movie"] = person.FavouriteMovie
		payload["favourite_color"] = person.FavouriteColor

		qdrant.UpsertPoints("people", vector, index, payload)

		fmt.Println(textToIndex)
	}
}

func fetchPeople() []Person {
	url := "https://tasks.aidevs.pl/data/people.json"
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Coult not read response")
	}

	articles := []Person{}
	err = json.Unmarshal(body, &articles)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	return articles
}
