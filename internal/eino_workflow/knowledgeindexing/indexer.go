package knowledgeindexing

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	redisCli "github.com/redis/go-redis/v9"
	"github.com/zhangga/aino/conf"
	"github.com/zhangga/aino/pkg/redispkg"
	"github.com/zhangga/aino/pkg/utils"
)

// newIndexer component initialization function of node 'Indexer1' in graph 'KnowledgeIndexing'
func newIndexer(ctx context.Context) (idr indexer.Indexer, err error) {
	if err = redispkg.Init(ctx, redispkg.Config{
		RedisAddr: conf.GlobalConfig.IndexerConf.RedisAddr,
		RedisPwd:  conf.GlobalConfig.IndexerConf.RedisPwd,
		Dimension: conf.GlobalConfig.IndexerConf.Dimension,
		Protocol:  conf.GlobalConfig.IndexerConf.Protocol,
	}, conf.GlobalConfig.IndexerConf.RedisPrefix); err != nil {
		return nil, err
	}

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

	// TODO Modify component configuration here.
	config := &redis.IndexerConfig{
		Client:    redisClient,
		KeyPrefix: conf.GlobalConfig.IndexerConf.RedisPrefix,
		BatchSize: 1,
		DocumentToHashes: func(ctx context.Context, doc *schema.Document) (*redis.Hashes, error) {
			if doc.ID == "" {
				//doc.ID = uuid.New().String()
				// 使用MD5摘要作为ID，保证相同内容只存一份
				sum := md5.Sum(utils.StringToBytes(doc.Content))
				doc.ID = hex.EncodeToString(sum[:])
			}
			key := doc.ID

			metadataBytes, err := json.Marshal(doc.MetaData)
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

	embeddingIns11, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns11
	idr, err = redis.NewIndexer(ctx, config)
	if err != nil {
		return nil, err
	}
	return idr, nil
}
