package zlog

var DefaultConfig = Config{
	FilePath: "./logs/app.log", // 默认日志文件路径
}

type Config struct {
	FilePath string `mapstructure:"file" json:"file" yaml:"file"`    // 日志文件路径
	Level    string `mapstructure:"level" json:"level" yaml:"level"` // 日志级别
}
