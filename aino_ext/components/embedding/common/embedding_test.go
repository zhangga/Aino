package common_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/zhangga/aino/aino_ext/components/embedding/common"
)

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
		"你好",
	}

	embeddings, err := embedder.EmbedStrings(ctx, texts)
	if err != nil {
		panic(err)
	}

	// 使用生成的向量
	for i, embedding := range embeddings {
		println("文本", i+1, "的向量维度:", len(embedding))
	}
}

func TestCommonEmbedding(t *testing.T) {
	ctx := context.Background()

	embedder, err := common.NewEmbedder(ctx, &common.EmbeddingConfig{
		BaseURL: os.Getenv("ARK_BASE_URL"),
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("ARK_EMBEDDING_MODEL"),
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
		println("文本", i+1, "的向量维度:", len(embedding))
	}
}
