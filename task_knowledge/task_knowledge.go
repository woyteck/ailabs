package task_knowledge

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"woyteck/ailabs/ailabs"
	"woyteck/ailabs/openapi"
)

type KnowledgeTask struct {
	Code      int64  `json:"code"`
	Msg       string `json:"msg"`
	Question  string `json:"question"`
	Database1 string `json:"database1"`
	Database2 string `json:"database2"`
}

type NBPRate struct {
	No            string  `json:"no"`
	EffectiveDate string  `json:"effectiveDate"`
	Mid           float64 `json:"mid"`
}

type NBPResponse struct {
	Table    string    `json:"table"`
	Currency string    `json:"currency"`
	Code     string    `json:"code"`
	Rates    []NBPRate `json:"rates"`
}

type CountryPopulationResponse struct {
	Population int `json:"population"`
}

func TaskKnowledge(apiKey string) {
	taskName := "knowledge"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result KnowledgeTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	detected, _ := detect(result.Question)
	fmt.Println(detected)
	if detected == "currency" {
		currency, _ := extractCurrency(result.Question)
		fmt.Println(currency)
		exchangeRate, _ := fetchExchangeRate(currency)
		fmt.Println(exchangeRate)
		answerResult := ailabs.AnswerAny(token, exchangeRate)
		fmt.Println(answerResult)
	} else if detected == "country population" {
		country, _ := extractCountry(result.Question)
		fmt.Println(country)
		population, _ := fetchCountryPopulation(country)
		fmt.Println(population)
		answerResult := ailabs.AnswerAny(token, population)
		fmt.Println(answerResult)
	} else {
		answer, _ := answerGeneralKnowledge(result.Question)
		fmt.Println(answer)
		answerResult := ailabs.AnswerString(token, answer)
		fmt.Println(answerResult)
	}
}

func fetchExchangeRate(currency string) (float64, error) {
	url := fmt.Sprintf("https://api.nbp.pl/api/exchangerates/rates/a/%v?format=json", currency)

	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return 0, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result NBPResponse
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return 0, err
	}

	return result.Rates[0].Mid, nil
}

func fetchCountryPopulation(country string) (int, error) {
	url := fmt.Sprintf("https://restcountries.com/v3.1/alpha/%v?fields=population", country)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result CountryPopulationResponse
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return 0, err
	}

	return result.Population, nil
}

func extractCurrency(text string) (string, error) {
	context := "I extract only the currency symbol and return it i 3-character ISO format, nothing else"

	messages := []openapi.Message{
		{Role: "system", Content: context},
		{Role: "user", Content: text},
	}
	response := openapi.GetCompletionShort(messages, "gpt-4")
	if len(response.Choices) == 0 {
		return "", errors.New("no completions returned")
	}

	return response.Choices[0].Message.Content, nil
}

func extractCountry(text string) (string, error) {
	context := "I extract the country symbol in ISO 3166-1 A-2 format and return it"

	messages := []openapi.Message{
		{Role: "system", Content: context},
		{Role: "user", Content: text},
	}
	response := openapi.GetCompletionShort(messages, "gpt-4")
	if len(response.Choices) == 0 {
		return "", errors.New("no completions returned")
	}

	return response.Choices[0].Message.Content, nil
}

func answerGeneralKnowledge(text string) (string, error) {
	context := "I return general knowledge answers. I keep it concise."

	messages := []openapi.Message{
		{Role: "system", Content: context},
		{Role: "user", Content: text},
	}
	response := openapi.GetCompletionShort(messages, "gpt-4")
	if len(response.Choices) == 0 {
		return "", errors.New("no completions returned")
	}

	return response.Choices[0].Message.Content, nil
}

func detect(text string) (string, error) {
	context := "I classify text to categories: currency, country population, other. I only return the category, nothing else"

	messages := []openapi.Message{
		{Role: "system", Content: context},
		{Role: "user", Content: text},
	}
	response := openapi.GetCompletionShort(messages, "gpt-4")
	if len(response.Choices) == 0 {
		return "", errors.New("no completions returned")
	}

	return response.Choices[0].Message.Content, nil
}
