package retriever_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

func init() {
	_ = godotenv.Load("../../../.env")
	fmt.Println(os.Getenv("ENV"))
}

func TestMilvusRetriever(t *testing.T) {
	addr := os.Getenv("INDEXER_MILVUS_ADDR")
	db := os.Getenv("INDEXER_MILVUS_DB")
	username := os.Getenv("INDEXER_MILVUS_USER")
	password := os.Getenv("INDEXER_MILVUS_PWD")
	arkApiKey := os.Getenv("ARK_API_KEY")
	arkModel := os.Getenv("ARK_EMBEDDING_MODEL")

	// Create a client
	ctx := context.Background()
	cli, err := client.NewClient(ctx, client.Config{
		Address:  addr,
		DBName:   db,
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
		Collection:  "test",
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
	documents, err := retriever.Retrieve(ctx, "太阳")
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
