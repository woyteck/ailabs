package indexer

import (
	"context"
	"flag"
	"time"

	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:6334", "the address to connect to")
)

var connection *grpc.ClientConn
var collectionsClient pb.CollectionsClient
var pointsClient pb.PointsClient

func getConnection() *grpc.ClientConn {
	if connection != nil {
		return connection
	}

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	// defer conn.Close()

	connection = conn

	return conn
}

func getCollectionsClient() pb.CollectionsClient {
	if collectionsClient == nil {
		collectionsClient = pb.NewCollectionsClient(getConnection())
	}

	return collectionsClient
}

func getPointsClient() pb.PointsClient {
	if pointsClient == nil {
		pointsClient = pb.NewPointsClient(connection)
	}

	return pointsClient
}

func CreateCollection(collectionName string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var defaultSegmentNumber uint64 = 2
	_, err := getCollectionsClient().Create(ctx, &pb.CreateCollection{
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

func Upsert(collectionName string, data []float64, num int, payloadStrings map[string]string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	payload := map[string]*pb.Value{}
	for index, value := range payloadStrings {
		payload[index] = &pb.Value{
			Kind: &pb.Value_StringValue{StringValue: value},
		}
	}

	embedding := []float32{}
	for _, f := range data {
		embedding = append(embedding, float32(f))
	}

	waitUpsert := true
	upsertPoints := []*pb.PointStruct{
		{
			Id: &pb.PointId{
				PointIdOptions: &pb.PointId_Num{Num: uint64(num)},
			},
			Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: embedding}}},
			Payload: payload,
		},
	}

	_, err := getPointsClient().Upsert(ctx, &pb.UpsertPoints{
		CollectionName: collectionName,
		Wait:           &waitUpsert,
		Points:         upsertPoints,
	})
	if err != nil {
		panic(err)
	}
}
