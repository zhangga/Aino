package common

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/zhangga/aino/pkg/httppkg"
)

var _ embedding.Embedder = (*EmbeddingImpl)(nil)

type EmbeddingImpl struct {
	config *EmbeddingConfig
}

type EmbeddingConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

type EmbeddingRequest struct {
	Model string      `json:"model"`
	Input []InputItem `json:"input"`
}

type InputItem struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type EmbeddingResponse struct {
	Created int64 `json:"created"`
	Data    struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
	Id     string `json:"id"`
	Model  string `json:"model"`
	Object string `json:"object"`
	Usage  struct {
		PromptTokens        int            `json:"prompt_tokens"`
		TotalTokens         int            `json:"total_tokens"`
		PromptTokensDetails map[string]int `json:"prompt_tokens_details"`
	} `json:"usage"`
}

func NewEmbedder(ctx context.Context, config *EmbeddingConfig) (emb embedding.Embedder, err error) {
	// TODO Modify component configuration here.
	emb = &EmbeddingImpl{config: config}
	return emb, nil
}

func (impl *EmbeddingImpl) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float64, error) {
	client := httppkg.NewDefaultClient(15 * time.Second)
	client.SetHeader("Authorization", "Bearer "+impl.config.APIKey)
	var inputs []InputItem
	for _, text := range texts {
		inputs = append(inputs, InputItem{
			Type: "text",
			Text: text,
		})
	}
	req := &httppkg.Request{
		Method: "POST",
		URL:    impl.config.BaseURL,
		Body: EmbeddingRequest{
			Model: impl.config.Model,
			Input: inputs,
		},
	}

	resp, err := client.Do(ctx, req)
	if err != nil {
		logger.Errorf("failed to get embedding: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("failed to get embedding, status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("failed to get embedding, status code: %d", resp.StatusCode)
	}

	var embeddingResp EmbeddingResponse
	if err = sonic.Unmarshal(resp.Body, &embeddingResp); err != nil {
		logger.Errorf("failed to unmarshal embedding response: %v", err)
		return nil, err
	}

	var results [][]float64
	results = append(results, embeddingResp.Data.Embedding)
	return results, nil
}
