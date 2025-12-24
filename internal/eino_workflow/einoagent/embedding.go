package einoagent

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/zhangga/aino/conf"
)

func newEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	// TODO Modify component configuration here.
	config := &ark.EmbeddingConfig{
		Model:  conf.GlobalConfig.EmbedConfig.Model,
		APIKey: conf.GlobalConfig.EmbedConfig.APIKey,
	}
	eb, err = ark.NewEmbedder(ctx, config)
	if err != nil {
		return nil, err
	}
	return eb, nil
}
