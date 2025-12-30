package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zhangga/aino/cmd/aino"
	"github.com/zhangga/aino/cmd/einoagent"
	"github.com/zhangga/aino/cmd/knowledgeindexing"
	"github.com/zhangga/aino/conf"
	"github.com/zhangga/aino/pkg/version"
	logger "github.com/zhangga/aino/pkg/zlog"
)

var rootCmd = &cobra.Command{
	Use:     "aino",
	Short:   "run aino service",
	Version: version.Version,
}

func init() {
	rootCmd.AddCommand(version.Command())
	rootCmd.AddCommand(aino.CmdRun)
	rootCmd.AddCommand(knowledgeindexing.CmdRun)
	rootCmd.AddCommand(einoagent.CmdRun)
}

var (
	ConfigPath string
)

// 需要绑定命令行参数的在这里注册
func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigPath, "config", "c", "configs/config.yaml", "config file path")

	// 日志配置
	defaultLogCfg := logger.DefaultConfig
	conf.GlobalConfig.LogConf = &defaultLogCfg
	rootCmd.Flags().StringVar(&conf.GlobalConfig.LogConf.FilePath, "log.file", "./logs/app.log", "log file path")
	rootCmd.Flags().StringVar(&conf.GlobalConfig.LogConf.Level, "log.level", "debug", "log level, eg: ENV: LOG_LEVEL=info")
	// 服务配置
	conf.GlobalConfig.ServiceConf = &conf.ServiceConfig{}
	rootCmd.Flags().IntVar(&conf.GlobalConfig.ServiceConf.HttpPort, "service.http_port", 8080, "service http port, eg: --service.http_port=8080")
	rootCmd.Flags().BoolVar(&conf.GlobalConfig.ServiceConf.Debug, "service.debug", true, "debug mode, eg: --service.debug=true")
	rootCmd.Flags().BoolVar(&conf.GlobalConfig.ServiceConf.EinoDebug, "service.eino_debug", false, "eino debug mode, eg: --service.eino_debug=true")
	rootCmd.Flags().BoolVar(&conf.GlobalConfig.ServiceConf.StreamMode, "service.stream_mode", false, "stream mode, eg: --service.stream_mode=true")
	rootCmd.Flags().StringVar(&conf.GlobalConfig.ServiceConf.APMPlusAppKey, "service.apmplus_app_key", "", "APMPlus App Key, eg: --service.apmplus_app_key=xxxxx")
	rootCmd.Flags().StringVar(&conf.GlobalConfig.ServiceConf.APMPlusRegion, "service.apmplus_region", "cn-beijing", "APMPlus Region, eg: --service.apmplus_region=cn-beijing")
	rootCmd.Flags().StringVar(&conf.GlobalConfig.ServiceConf.LangfusePublicKey, "service.langfuse_public_key", "", "Langfuse Public Key, eg: --service.langfuse_public_key=xxxxx")
	rootCmd.Flags().StringVar(&conf.GlobalConfig.ServiceConf.LangfuseSecretKey, "service.langfuse_secret_key", "", "Langfuse Secret Key, eg: --service.langfuse_secret_key=xxxxx")
	// Embed配置
	conf.GlobalConfig.EmbedConf = &conf.EmbedConfig{}
	rootCmd.Flags().StringVar(&conf.GlobalConfig.EmbedConf.BaseURL, "embed.base_url", "", "embedding url, eg: --embed.base_url=https://ark.cn-beijing.volces.com/api/v3")
	rootCmd.Flags().StringVar(&conf.GlobalConfig.EmbedConf.APIKey, "embed.api_key", "", "embedding api key, eg: --embed.api_key=xxxxx")
	rootCmd.Flags().StringVar(&conf.GlobalConfig.EmbedConf.Model, "embed.model", "", "embedding model, eg: --embed.model=xxxxx")
	// Indexer配置
	conf.GlobalConfig.IndexerConf = &conf.IndexerConfig{}
	rootCmd.Flags().IntVar(&conf.GlobalConfig.IndexerConf.Dimension, "indexer.dimension", 2048, "vector dimension, eg: --indexer.dimension=2048")
	conf.GlobalConfig.IndexerConf.RedisConf = &conf.RedisConfig{}
	rootCmd.Flags().StringVar(&conf.GlobalConfig.IndexerConf.RedisConf.IndexPrefix, "indexer.redis.index_prefix", "aino_doc:", "redis index prefix, eg: --indexer.redis.index_prefix=aino_doc:")
	rootCmd.Flags().StringVar(&conf.GlobalConfig.IndexerConf.RedisConf.IndexName, "indexer.redis.index_name", "vector_index", "redis index name, eg: --indexer.redis.index_name=vector_index")
	rootCmd.Flags().IntVar(&conf.GlobalConfig.IndexerConf.RedisConf.Protocol, "indexer.redis.protocol", 2, "redis protocol, eg: --indexer.redis.protocol=2")
	conf.GlobalConfig.IndexerConf.MilvusConf = &conf.MilvusConfig{}
	rootCmd.Flags().StringVar(&conf.GlobalConfig.IndexerConf.MilvusConf.DBName, "indexer.milvus.db_name", "aino", "milvus dbName, eg: --indexer.milvus.db_name=aino")
	rootCmd.Flags().StringVar(&conf.GlobalConfig.IndexerConf.MilvusConf.Collection, "indexer.milvus.collection", "doc", "milvus collection name, eg: --indexer.milvus.collection=doc")

}

func checkConfigPath() {
	lastDot := strings.LastIndex(ConfigPath, ".")
	if lastDot == -1 {
		return
	}

	ext := ConfigPath[lastDot+1:]
	// 有需要的话可以根据环境变换下配置文件的路径
	NewConfigPath := fmt.Sprintf("%s.%s", ConfigPath[:lastDot], ext)
	// 检查文件是否存在
	if _, err := os.Stat(NewConfigPath); err == nil {
		ConfigPath = NewConfigPath
	}
}

func init() {
	// 配置优先级：命令行参数 > 环境变量 >.env > 配置文件
	checkConfigPath()
	viper.SetConfigFile(ConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config file: --config=%s failed, err=%s", ConfigPath, err)
	}

	// 2. 环境变量（中间优先级）
	_ = godotenv.Load()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// 3. 命令行参数（最高优先级，只绑定实际传入的）
	rootCmd.Flags().Visit(func(f *pflag.Flag) {
		_ = viper.BindPFlag(f.Name, f)
	})

	// 4. 最终将配置文件内容解析到Config结构体中
	if err := viper.Unmarshal(&conf.GlobalConfig); err != nil {
		log.Fatalf("unmarshal config file failed, err=%s", err)
	}

	// 5. 检查必要的配置项
	if !conf.GlobalConfig.ValidData() {
		log.Fatalf("config validation failed: %v", conf.GlobalConfig)
	}
	log.Printf("load config from %s success", ConfigPath)
}
