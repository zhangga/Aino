package knowledgeindexing

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/components/embedding"
)

type EmbeddingImpl struct {
	config *EmbeddingConfig
}

type EmbeddingConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

func newEmbedding(ctx context.Context) (emb embedding.Embedder, err error) {
	// TODO Modify component configuration here.
	config := &EmbeddingConfig{
		BaseURL: "https://api.example.com/embeddings",
		APIKey:  "your_api_key",
		Model:   "example-embedding-model",
	}
	emb = &EmbeddingImpl{config: config}
	return emb, nil
}

func (impl *EmbeddingImpl) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float64, error) {
	eb, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		BaseURL: impl.config.BaseURL,
		APIKey:  impl.config.APIKey,
		Model:   impl.config.Model,
	})
	if err != nil {
		return nil, err
	}
	return eb.EmbedStrings(ctx, texts, opts...)
}
