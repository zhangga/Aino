package knowledgeindexing

import (
	"context"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/zhangga/aino/aino_ext/components/embedding/common"
	"github.com/zhangga/aino/conf"
)

func newEmbedding(ctx context.Context) (emb embedding.Embedder, err error) {
	return common.NewEmbedder(ctx, &common.EmbeddingConfig{
		BaseURL: conf.GlobalConfig.EmbedConfig.BaseURL,
		APIKey:  conf.GlobalConfig.EmbedConfig.APIKey,
		Model:   conf.GlobalConfig.EmbedConfig.Model,
	})
}
