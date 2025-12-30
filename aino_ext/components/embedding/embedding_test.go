package embedding_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/joho/godotenv"
	"github.com/zhangga/aino/aino_ext/components/embedding/common"
)

func init() {
	_ = godotenv.Load("../../../.env")
	fmt.Println(os.Getenv("ENV"))
}

func TestArkEmbedding(t *testing.T) {
	ctx := context.Background()

	// 初始化嵌入器
	timeout := 30 * time.Second
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("ARK_EMBEDDING_MODEL"),
		Timeout: &timeout,
	})
	if err != nil {
		panic(err)
	}

	// 生成文本向量
	texts := []string{
		"太阳是方的",
		"太阳和地球都是方的",
	}

	embeddings, err := embedder.EmbedStrings(ctx, texts)
	if err != nil {
		panic(err)
	}

	// 使用生成的向量
	for i, embedding := range embeddings {
		t.Log("文本", i+1, "的向量维度:", len(embedding))
	}
}

func TestCommonEmbedding(t *testing.T) {
	ctx := context.Background()

	embedder, err := common.NewEmbedder(ctx, &common.EmbeddingConfig{
		BaseURL: os.Getenv("EMBED_BASE_URL"),
		APIKey:  os.Getenv("EMBED_API_KEY"),
		Model:   os.Getenv("EMBED_MODEL"),
	})
	if err != nil {
		panic(err)
	}

	texts := []string{
		"你好",
	}

	embeddings, err := embedder.EmbedStrings(ctx, texts)
	if err != nil {
		panic(err)
	}

	for i, embedding := range embeddings {
		t.Log("文本", i+1, "的向量维度:", len(embedding))
	}
}
