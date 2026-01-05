package retriever_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	redisCli "github.com/redis/go-redis/v9"
	"github.com/zhangga/aino/pkg/redispkg"
)

func init() {
	_ = godotenv.Load("../../../.env")
	fmt.Println(os.Getenv("ENV"))
}

func TestMilvusRetriever(t *testing.T) {
	addr := os.Getenv("INDEXER_MILVUS_ADDR")
	username := os.Getenv("INDEXER_MILVUS_USER")
	password := os.Getenv("INDEXER_MILVUS_PWD")
	dbName := "indexer_test"
	collectionName := "test"
	arkApiKey := os.Getenv("ARK_API_KEY")
	arkModel := os.Getenv("ARK_EMBEDDING_MODEL")

	// Create a client
	ctx := context.Background()
	cli, err := client.NewClient(ctx, client.Config{
		Address:  addr,
		DBName:   dbName,
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
		return
	}
	defer cli.Close()

	// Create an embedding model
	emb, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: arkApiKey,
		Model:  arkModel,
	})

	// Create a retriever
	retriever, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:      cli,
		Collection:  collectionName,
		Partition:   nil,
		VectorField: "",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		DocumentConverter: nil,
		MetricType:        "",
		TopK:              2,
		ScoreThreshold:    5,
		Sp:                nil,
		Embedding:         emb,
	})
	if err != nil {
		t.Fatalf("Failed to create retriever: %v", err)
		return
	}

	// Retrieve documents
	documents, err := retriever.Retrieve(ctx, "张泽强")
	if err != nil {
		t.Fatalf("Failed to retrieve: %v", err)
		return
	}

	// Print the documents
	for i, doc := range documents {
		t.Logf("Document %d:", i)
		t.Logf("title: %s\n", doc.ID)
		t.Logf("content: %s\n", doc.Content)
		t.Logf("metadata: %v\n", doc.MetaData)
	}
}

func TestRedisRetriever(t *testing.T) {
	addr := os.Getenv("INDEXER_REDIS_ADDR")
	password := os.Getenv("INDEXER_REDIS_PWD")
	indexPrefix := "indexer_test:"
	indexName := "vector_index"
	arkApiKey := os.Getenv("ARK_API_KEY")
	arkModel := os.Getenv("ARK_EMBEDDING_MODEL")

	ctx := context.Background()

	var redisClient *redisCli.Client
	if len(password) == 0 {
		redisClient = redisCli.NewClient(&redisCli.Options{
			Addr:     addr,
			Protocol: 2,
		})
	} else {
		redisClient = redisCli.NewClient(&redisCli.Options{
			Addr:     addr,
			Password: password,
			Protocol: 2,
		})
	}

	config := &redis.RetrieverConfig{
		Client:       redisClient,
		Index:        fmt.Sprintf("%s%s", indexPrefix, indexName),
		Dialect:      2,
		ReturnFields: []string{redispkg.ContentField, redispkg.MetadataField, redispkg.DistanceField},
		TopK:         8,
		VectorField:  redispkg.VectorField,
		DocumentConverter: func(ctx context.Context, doc redisCli.Document) (*schema.Document, error) {
			resp := &schema.Document{
				ID:       doc.ID,
				Content:  "",
				MetaData: map[string]any{},
			}
			for field, val := range doc.Fields {
				if field == redispkg.ContentField {
					resp.Content = val
				} else if field == redispkg.MetadataField {
					resp.MetaData[field] = val
				} else if field == redispkg.DistanceField {
					distance, err := strconv.ParseFloat(val, 64)
					if err != nil {
						continue
					}
					resp.WithScore(1 - distance)
				}
			}

			return resp, nil
		},
	}

	// Create an embedding model
	emb, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: arkApiKey,
		Model:  arkModel,
	})
	if err != nil {
		t.Fatalf("Failed to create embedding: %v", err)
	}
	config.Embedding = emb

	rtr, err := redis.NewRetriever(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create retriever: %v", err)
	}

	// Retrieve documents
	documents, err := rtr.Retrieve(ctx, "张泽强", retriever.WithTopK(1))
	if err != nil {
		t.Fatalf("Failed to retrieve: %v", err)
		return
	}

	// Print the documents
	for i, doc := range documents {
		t.Logf("Document %d:", i)
		t.Logf("title: %s\n", doc.ID)
		t.Logf("content: %s\n", doc.Content)
		t.Logf("metadata: %v\n", doc.MetaData)
	}
}
