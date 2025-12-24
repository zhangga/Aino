package conf

import logger "github.com/zhangga/aino/pkg/zlog"

var GlobalConfig Config

type IConfig interface {
	ValidData() bool
}
type Config struct {
	LogConf     *logger.Config `mapstructure:"log" json:"log" yaml:"log"`       // 日志相关配置
	EmbedConfig *EmbedConfig   `mapstructure:"embed" json:"embed" yaml:"embed"` // 向量配置
	IndexerConf *IndexerConfig `mapstructure:"indexer" json:"indexer" yaml:"indexer"`
}

type EmbedConfig struct {
	BaseURL string `mapstructure:"base_url" json:"base_url" yaml:"base_url"` // 向量服务地址
	APIKey  string `mapstructure:"api_key" json:"api_key" yaml:"api_key"`    // 向量服务API Key
	Model   string `mapstructure:"model" json:"model" yaml:"model"`          // 向量模型
}

type IndexerConfig struct {
	RedisAddr   string `mapstructure:"redis_addr" json:"redis_addr" yaml:"redis_addr"`       // Redis地址
	RedisPwd    string `mapstructure:"redis_pwd" json:"redis_pwd" yaml:"redis_pwd"`          // Redis密码
	RedisPrefix string `mapstructure:"redis_prefix" json:"redis_prefix" yaml:"redis_prefix"` // Redis键前缀
	Dimension   int    `mapstructure:"dimension" json:"dimension" yaml:"dimension"`          // 向量维度
	Protocol    int    `mapstructure:"protocol" json:"protocol" yaml:"protocol"`             // Redis协议版本
}

func (c *Config) ValidData() bool {
	return true
}
