package indexer_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func init() {
	_ = godotenv.Load("../../../.env")
	fmt.Println(os.Getenv("ENV"))
}

func TestMilvusIndexer(t *testing.T) {
	ctx := context.Background()

	milvusClient, err := client.NewClient(ctx, client.Config{
		Address: "localhost:19530",
		DBName:  "test_indexer",
	})
	if err != nil {
		t.Fatal(err)
	}

	fields := []*entity.Field{
		{
			Name:     "id",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "255",
			},
			PrimaryKey: true,
		},
		{
			Name:     "vector",
			DataType: entity.FieldTypeBinaryVector,
			TypeParams: map[string]string{
				"dim": "65536",
			},
		},
		{
			Name:     "content",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "8192",
			},
		},
		{
			Name:     "metadata",
			DataType: entity.FieldTypeJSON,
		},
	}

	//embedder, err := common.NewEmbedder(ctx, &common.EmbeddingConfig{
	//	BaseURL: os.Getenv("EMBED_BASE_URL"),
	//	APIKey:  os.Getenv("EMBED_API_KEY"),
	//	Model:   os.Getenv("EMBED_MODEL"),
	//})
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("ARK_EMBEDDING_MODEL"),
	})

	// 初始化索引器
	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:     milvusClient,
		Collection: "test",
		Fields:     fields,
		Embedding:  embedder,
	})
	if err != nil {
		t.Fatalf("failed to create Milvus indexer: %v", err)
	}

	docs := []*schema.Document{
		{
			ID:      "1",
			Content: "太阳是方的",
			MetaData: map[string]any{
				"author": "Alice",
			},
		},
		{
			ID:      "2",
			Content: "太阳和地球都是方的",
			MetaData: map[string]any{
				"author": "Bob",
			},
		},
	}

	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		t.Fatalf("failed to store documents: %v", err)
	}
	t.Logf("stored documents: %v", ids)
}
