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
		Dimension: conf.GlobalConfig.IndexerConf.Dimension,
		RedisAddr: conf.GlobalConfig.IndexerConf.RedisConf.Addr,
		RedisPwd:  conf.GlobalConfig.IndexerConf.RedisConf.Pwd,
		IndexName: conf.GlobalConfig.IndexerConf.RedisConf.IndexName,
		Protocol:  conf.GlobalConfig.IndexerConf.RedisConf.Protocol,
	}, conf.GlobalConfig.IndexerConf.RedisConf.IndexPrefix); err != nil {
		return nil, err
	}

	var redisClient *redisCli.Client
	if len(conf.GlobalConfig.IndexerConf.RedisConf.Pwd) == 0 {
		redisClient = redisCli.NewClient(&redisCli.Options{
			Addr:     conf.GlobalConfig.IndexerConf.RedisConf.Addr,
			Protocol: conf.GlobalConfig.IndexerConf.RedisConf.Protocol,
		})
	} else {
		redisClient = redisCli.NewClient(&redisCli.Options{
			Addr:     conf.GlobalConfig.IndexerConf.RedisConf.Addr,
			Password: conf.GlobalConfig.IndexerConf.RedisConf.Pwd,
			Protocol: conf.GlobalConfig.IndexerConf.RedisConf.Protocol,
		})
	}

	// TODO Modify component configuration here.
	config := &redis.IndexerConfig{
		Client:    redisClient,
		KeyPrefix: conf.GlobalConfig.IndexerConf.RedisConf.IndexPrefix,
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
