package einoagent

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	redisCli "github.com/redis/go-redis/v9"
	"github.com/zhangga/aino/conf"
	"github.com/zhangga/aino/pkg/redispkg"
)

// newRetriever component initialization function of node 'Retriever1' in graph 'EinoAgent'
func newRetriever(ctx context.Context) (rtr retriever.Retriever, err error) {
	// TODO Modify component configuration here.
	var redisClient *redisCli.Client
	if len(conf.GlobalConfig.IndexerConf.RedisPwd) == 0 {
		redisClient = redisCli.NewClient(&redisCli.Options{
			Addr:     conf.GlobalConfig.IndexerConf.RedisAddr,
			Protocol: conf.GlobalConfig.IndexerConf.Protocol,
		})
	} else {
		redisClient = redisCli.NewClient(&redisCli.Options{
			Addr:     conf.GlobalConfig.IndexerConf.RedisAddr,
			Password: conf.GlobalConfig.IndexerConf.RedisPwd,
			Protocol: conf.GlobalConfig.IndexerConf.Protocol,
		})
	}
	config := &redis.RetrieverConfig{
		Client:       redisClient,
		Index:        fmt.Sprintf("%s%s", conf.GlobalConfig.IndexerConf.RedisPrefix, conf.GlobalConfig.IndexerConf.IndexName),
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
	embeddingIns11, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns11
	rtr, err = redis.NewRetriever(ctx, config)
	if err != nil {
		return nil, err
	}
	return rtr, nil
}
