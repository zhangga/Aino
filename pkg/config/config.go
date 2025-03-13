package config

import (
	"errors"
	"fmt"
	"github.com/zhangga/aino/pkg/logger"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type IConfig interface {
	ValidData() bool
}

var (
	ErrorConfigPathEmpty   = errors.New("config path is empty")
	ErrorConfigFileNoExist = errors.New("config file no exist")
)

// LoadConfig 读取项目配置参数
// 命令行 > 环境变量 > .env > 配置文件
func LoadConfig(cc IConfig, configPath string) error {
	var err error
	// 先加载.env文件
	_ = godotenv.Load()

	// 从配置文件读取
	if err = loadConfigByFile(configPath); err != nil {
		return err
	}

	// 从环境变量中读取
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 反序列化配置
	if err = viper.Unmarshal(cc); err != nil {
		return err
	}

	// 检查必要的参数是否设置
	if !cc.ValidData() {
		return fmt.Errorf("config check failed, configPath=%s", configPath)
	}
	return nil
}

// loadConfigByFile 从配置文件中读取
func loadConfigByFile(configPath string) error {
	// 未指定配置文件
	if len(configPath) == 0 {
		return ErrorConfigPathEmpty
	}

	var fileName, fileSuffix string
	// 获取文件名和后缀
	if idx := strings.LastIndex(configPath, "."); idx >= 0 {
		fileName = configPath[:idx]
		fileSuffix = configPath[idx+1:]
	} else {
		return fmt.Errorf("----->: configPath=%s, must has file suffix", configPath)
	}

	// 1. 读取对应env的配置文件. 如: configs/config_dev.toml
	flagEnv := os.Getenv("ENV")
	filePath := fmt.Sprintf("%s_%s.%s", fileName, flagEnv, fileSuffix)
	_, err := os.Stat(filePath)
	// 文件存在
	if err == nil || os.IsExist(err) {
		logger.Infof("load config file ----->: %s", filePath)
		viper.SetConfigFile(filePath)
		return viper.ReadInConfig()
	}

	// 2. 读取指定的配置文件
	_, err = os.Stat(configPath)
	if err == nil || os.IsExist(err) {
		logger.Infof("load config file ----->: %s", configPath)
		viper.SetConfigFile(configPath)
		return viper.ReadInConfig()
	}

	// 3. 没找到配置文件
	return ErrorConfigFileNoExist
}
