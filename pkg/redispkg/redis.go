package redispkg

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
)

const (
	defaultRedisPrefix = "knowledge:doc:"
	defaultIndexName   = "vector_index"

	ContentField  = "content"
	MetadataField = "metadata"
	VectorField   = "content_vector"
	DistanceField = "distance"
)

var initOnce sync.Once

func Init(ctx context.Context, config Config, redisPrefix string) error {
	var err error
	initOnce.Do(func() {
		err = InitRedisIndex(ctx, &config, redisPrefix)
	})
	return err
}

type Config struct {
	RedisAddr string
	RedisPwd  string
	IndexName string
	Dimension int
	Protocol  int
}

func InitRedisIndex(ctx context.Context, config *Config, redisPrefix string) (err error) {
	if len(redisPrefix) == 0 {
		redisPrefix = defaultRedisPrefix // 默认前缀
	}
	indexName := config.IndexName
	if len(indexName) == 0 {
		indexName = defaultIndexName // 默认索引名
	}
	if config.Dimension <= 0 {
		return fmt.Errorf("dimension must be positive")
	}

	var client *redis.Client
	if len(config.RedisPwd) == 0 {
		client = redis.NewClient(&redis.Options{
			Addr:     config.RedisAddr,
			Protocol: config.Protocol,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     config.RedisAddr,
			Password: config.RedisPwd,
			Protocol: config.Protocol,
		})
	}

	defer func() {
		if err != nil {
			_ = client.Close()
		}
	}()

	if err = client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	indexNameFull := fmt.Sprintf("%s%s", redisPrefix, indexName)

	// 检查是否存在索引
	exists, err := client.Do(ctx, "FT.INFO", indexNameFull).Result()
	if err != nil {
		if !strings.Contains(err.Error(), "Unknown index name") {
			return fmt.Errorf("failed to check if index exists: %w", err)
		}
		err = nil
	} else if exists != nil {
		// 删除索引
		//drop, err := client.Do(ctx, "FT.DROPINDEX", indexNameFull).Result()
		//if err != nil {
		//	return err
		//}
		//_ = drop
		return nil
	}

	// Create new index
	createIndexArgs := []interface{}{
		"FT.CREATE", indexNameFull,
		"ON", "HASH",
		"PREFIX", "1", redisPrefix,
		"SCHEMA",
		ContentField, "TEXT",
		MetadataField, "TEXT",
		VectorField, "VECTOR", "FLAT",
		"6",
		"TYPE", "FLOAT32",
		"DIM", config.Dimension,
		"DISTANCE_METRIC", "COSINE",
	}

	if err = client.Do(ctx, createIndexArgs...).Err(); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// 验证索引是否创建成功
	if _, err = client.Do(ctx, "FT.INFO", indexNameFull).Result(); err != nil {
		return fmt.Errorf("failed to verify index creation: %w", err)
	}

	return nil
}
