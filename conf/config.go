package conf

import logger "github.com/zhangga/aino/pkg/zlog"

var GlobalConfig Config

type IConfig interface {
	ValidData() bool
}
type Config struct {
	LogConf     *logger.Config `mapstructure:"log" json:"log" yaml:"log"`             // 日志相关配置
	ServiceConf *ServiceConfig `mapstructure:"service" json:"service" yaml:"service"` // 服务相关配置
	LLMConf     *LLMConfig     `mapstructure:"llm" json:"llm" yaml:"llm"`             // 大模型配置
	EmbedConf   *EmbedConfig   `mapstructure:"embed" json:"embed" yaml:"embed"`       // 向量配置
	IndexerConf *IndexerConfig `mapstructure:"indexer" json:"indexer" yaml:"indexer"` // 索引配置
}

type ServiceConfig struct {
	HttpPort          int    `mapstructure:"http_port" json:"http_port" yaml:"http_port"`                               // 服务端口
	Debug             bool   `mapstructure:"debug" json:"debug" yaml:"debug"`                                           // 调试模式
	EinoDebug         bool   `mapstructure:"eino_debug" json:"eino_debug" yaml:"eino_debug"`                            // Eino调试模式
	StreamMode        bool   `mapstructure:"stream_mode" json:"stream_mode" yaml:"stream_mode"`                         // 流式模式
	APMPlusAppKey     string `mapstructure:"apmplus_app_key" json:"apmplus_app_key" yaml:"apmplus_app_key"`             // APMPlus应用Key
	APMPlusRegion     string `mapstructure:"apmplus_region" json:"apmplus_region" yaml:"apmplus_region"`                // APMPlus区域
	LangfusePublicKey string `mapstructure:"langfuse_public_key" json:"langfuse_public_key" yaml:"langfuse_public_key"` // Langf使用公钥
	LangfuseSecretKey string `mapstructure:"langfuse_secret_key" json:"langfuse_secret_key" yaml:"langfuse_secret_key"` // Langf使用私钥
}

type LLMConfig struct {
	BaseURL string `mapstructure:"base_url" json:"base_url" yaml:"base_url"` // 大模型服务地址
	Model   string `mapstructure:"model" json:"model" yaml:"model"`          // 大模型名称
	APIKey  string `mapstructure:"api_key" json:"api_key" yaml:"api_key"`    // 大模型API Key
}

type EmbedConfig struct {
	BaseURL string `mapstructure:"base_url" json:"base_url" yaml:"base_url"` // 向量服务地址
	APIKey  string `mapstructure:"api_key" json:"api_key" yaml:"api_key"`    // 向量服务API Key
	Model   string `mapstructure:"model" json:"model" yaml:"model"`          // 向量模型
}

type IndexerConfig struct {
	Dimension  int           `mapstructure:"dimension" json:"dimension" yaml:"dimension"` // 向量维度
	RedisConf  *RedisConfig  `mapstructure:"redis" json:"redis" yaml:"redis"`             // Redis配置
	MilvusConf *MilvusConfig `mapstructure:"milvus" json:"milvus" yaml:"milvus"`          // Milvus配置
}

type RedisConfig struct {
	Addr        string `mapstructure:"addr" json:"addr" yaml:"addr"`                         // Redis地址
	Pwd         string `mapstructure:"pwd" json:"pwd" yaml:"pwd"`                            // Redis密码
	IndexPrefix string `mapstructure:"index_prefix" json:"index_prefix" yaml:"index_prefix"` // 索引前缀
	IndexName   string `mapstructure:"index_name" json:"index_name" yaml:"index_name"`       // 索引名称
	Protocol    int    `mapstructure:"protocol" json:"protocol" yaml:"protocol"`             // Redis协议版本
}

type MilvusConfig struct {
	Addr       string `mapstructure:"addr" json:"addr" yaml:"addr"`                   // Milvus地址
	User       string `mapstructure:"user" json:"user" yaml:"user"`                   // Milvus用户
	Pwd        string `mapstructure:"pwd" json:"pwd" yaml:"pwd"`                      // Milvus密码
	DBName     string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`          // Milvus数据库名称
	Collection string `mapstructure:"collection" json:"collection" yaml:"collection"` // Milvus集合名称
}

func (c *Config) ValidData() bool {
	return true
}
