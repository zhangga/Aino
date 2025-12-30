package indexer_test

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	redisCli "github.com/redis/go-redis/v9"
	"github.com/zhangga/aino/pkg/redispkg"
	"github.com/zhangga/aino/pkg/utils"
)

func init() {
	_ = godotenv.Load("../../../.env")
	fmt.Println(os.Getenv("ENV"))
}

func TestMilvusIndexer(t *testing.T) {
	addr := os.Getenv("INDEXER_MILVUS_ADDR")
	username := os.Getenv("INDEXER_MILVUS_USER")
	password := os.Getenv("INDEXER_MILVUS_PWD")
	dbName := "indexer_test"
	collectionName := "test"
	arkApiKey := os.Getenv("ARK_API_KEY")
	arkModel := os.Getenv("ARK_EMBEDDING_MODEL")

	ctx := context.Background()

	milvusClient, err := client.NewClient(ctx, client.Config{
		Address:  addr,
		DBName:   dbName,
		Username: username,
		Password: password,
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
		APIKey: arkApiKey,
		Model:  arkModel,
	})

	// 初始化索引器
	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:     milvusClient,
		Collection: collectionName,
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

func TestRedisIndexer(t *testing.T) {
	addr := os.Getenv("INDEXER_REDIS_ADDR")
	password := os.Getenv("INDEXER_REDIS_PWD")
	indexPrefix := "indexer_test:"
	indexName := "vector_index"
	arkApiKey := os.Getenv("ARK_API_KEY")
	arkModel := os.Getenv("ARK_EMBEDDING_MODEL")

	ctx := context.Background()

	// 初始化 Redis 索引
	if err := redispkg.InitRedisIndex(ctx, &redispkg.Config{
		RedisAddr: addr,
		RedisPwd:  password,
		IndexName: indexName,
		Dimension: 2048,
		Protocol:  2,
	}, indexPrefix); err != nil {
		t.Fatalf("failed to init redis index: %v", err)
	}

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

	config := &redis.IndexerConfig{
		Client:    redisClient,
		KeyPrefix: indexPrefix,
		BatchSize: 1,
		DocumentToHashes: func(ctx context.Context, doc *schema.Document) (*redis.Hashes, error) {
			if doc.ID == "" {
				//doc.ID = uuid.New().String()
				// 使用MD5摘要作为ID，保证相同内容只存一份
				sum := md5.Sum(utils.StringToBytes(doc.Content))
				doc.ID = hex.EncodeToString(sum[:])
			}
			key := doc.ID

			metadataBytes, err := sonic.Marshal(doc.MetaData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal metadata: %w", err)
			}

			return &redis.Hashes{
				Key: key,
				Field2Value: map[string]redis.FieldValue{
					redispkg.ContentField:  {Value: doc.Content, EmbedKey: redispkg.VectorField},
					redispkg.MetadataField: {Value: metadataBytes},
				},
			}, nil
		},
	}

	//embedder, err := common.NewEmbedder(ctx, &common.EmbeddingConfig{
	//	BaseURL: os.Getenv("EMBED_BASE_URL"),
	//	APIKey:  os.Getenv("EMBED_API_KEY"),
	//	Model:   os.Getenv("EMBED_MODEL"),
	//})
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: arkApiKey,
		Model:  arkModel,
	})
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}

	config.Embedding = embedder
	indexer, err := redis.NewIndexer(ctx, config)
	if err != nil {
		t.Fatalf("failed to create Redis indexer: %v", err)
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
		{
			ID:      "3",
			Content: "张泽强是一个大帅哥",
			MetaData: map[string]any{
				"author": "Jossy",
			},
		},
	}

	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		t.Fatalf("failed to store documents: %v", err)
	}
	t.Logf("stored documents: %v", ids)
}
