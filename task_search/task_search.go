package task_search

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"woyteck/ailabs/ailabs"

	"github.com/google/uuid"
	pb "github.com/qdrant/go-client/qdrant"
	"github.com/tmc/langchaingo/llms/openai"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:6334", "the address to connect to")
)

var collectionName string = "ailabs"

type SearchTask struct {
	Code     int64  `json:"code"`
	Msg      string `json:"msg"`
	Question string `json:"question"`
}

type Article struct {
	UUID  uuid.UUID `json:"uuid"`
	Title string    `json:"title"`
	Url   string    `json:"url"`
	Info  string    `json:"info"`
	Date  string    `json:"date"`
}

func TaskSearch(apiKey string) {
	taskName := "search"
	token := ailabs.GetToken(taskName, apiKey)
	taskJsonString := ailabs.GetTask(token)
	fmt.Println(string(taskJsonString))

	var result SearchTask
	err := json.Unmarshal(taskJsonString, &result)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	llm := createLlm()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// create collection
	collectionsClient := pb.NewCollectionsClient(conn)
	createCollection(collectionsClient)

	// index data
	pointsClient := pb.NewPointsClient(conn)
	indexArticles(llm, pointsClient, getArticlesList())

	// search
	searchPhraseVector := embedString(llm, result.Question)
	results, _ := search(pointsClient, collectionName, searchPhraseVector)
	if len(results) == 0 {
		panic("no results")
	}

	//return answer
	payload := results[0].GetPayload()
	url := payload["url"].GetStringValue()

	answerResult := ailabs.AnswerString(token, url)
	fmt.Println(answerResult)
}

func search(pointsClient pb.PointsClient, collectionName string, vector []float32) ([]*pb.ScoredPoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	searchResult, err := pointsClient.Search(ctx, &pb.SearchPoints{
		CollectionName: collectionName,
		Vector:         vector,
		Limit:          1,
		WithPayload:    &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		return nil, err
	}

	return searchResult.GetResult(), err
}

func indexArticles(llm *openai.LLM, pointsClient pb.PointsClient, articles []Article) {

	for index, article := range articles {
		embedding := embedString(llm, article.Title)
		upsert(pointsClient, embedding, article.Url, index+1)
	}
}

func embedString(llm *openai.LLM, text string) []float32 {
	ctx := context.Background()
	embedings, err := llm.CreateEmbedding(ctx, []string{text})
	if err != nil {
		panic(err)
	}

	return embedings[0]
}

func upsert(pointsClient pb.PointsClient, data []float32, url string, num int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	waitUpsert := true
	upsertPoints := []*pb.PointStruct{
		{
			Id: &pb.PointId{
				PointIdOptions: &pb.PointId_Num{Num: uint64(num)},
			},
			Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: data}}},
			Payload: map[string]*pb.Value{
				"url": {
					Kind: &pb.Value_StringValue{StringValue: url},
				},
			},
		},
	}

	_, err := pointsClient.Upsert(ctx, &pb.UpsertPoints{
		CollectionName: collectionName,
		Wait:           &waitUpsert,
		Points:         upsertPoints,
	})
	if err != nil {
		panic(err)
	}

	log.Println("Upsert", len(upsertPoints), "points")
}

func createCollection(collectionsClient pb.CollectionsClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var defaultSegmentNumber uint64 = 2
	_, err := collectionsClient.Create(ctx, &pb.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: &pb.VectorsConfig{
			Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     1536,
					Distance: pb.Distance_Cosine,
				},
			},
		},
		OptimizersConfig: &pb.OptimizersConfigDiff{
			DefaultSegmentNumber: &defaultSegmentNumber,
		},
	})
	if err != nil {
		panic(err)
	}
}

func createLlm() *openai.LLM {
	opts := []openai.Option{
		openai.WithModel("gpt-3.5-turbo-0125"),
		openai.WithEmbeddingModel("text-embedding-ada-002"),
	}
	llm, err := openai.New(opts...)
	if err != nil {
		log.Fatal(err)
	}

	return llm
}

func getArticlesList() []Article {
	url := "https://unknow.news/archiwum_aidevs.json"
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error occured %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Coult not read response")
	}

	articles := []Article{}
	err = json.Unmarshal(body, &articles)
	if err != nil {
		log.Fatal("Can not unmarshall JSON")
	}

	return articles
}
