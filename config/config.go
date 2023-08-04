package config

import (
	"os"
	"regexp"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Jaeger   JaegerConfig   `mapstructure:"jaeger"`
}

type ServerConfig struct {
	Host           string        `mapstructure:"host"`
	Port           string        `mapstructure:"port"`
	Debug          bool          `mapstructure:"debug"`
	ReadTimeout    time.Duration `mapstructure:"readTimeout"`
	WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
	ContextTimeout time.Duration `mapstructure:"contextTimeout"`
	Timezone       string        `mapstructure:"timezone"`
}

type DatabaseConfig struct {
	Adapter         string `mapstructure:"adapter"`
	Host            string `mapstructure:"host"`
	Username        string `mapstructure:"username"`
	Db              string `mapstructure:"db"`
	Password        string `mapstructure:"password"`
	Port            int    `mapstructure:"port"`
	MaxConns        int    `mapstructure:"maxConns"`
	MaxLiftimeConns int    `mapstructure:"maxLiftimeConns"`
}

type LoggerConfig struct {
	Development       bool   `mapstructure:"development"`
	DisableCaller     bool   `mapstructure:"disableCaller"`
	DisableStacktrace bool   `mapstructure:"disableStacktrace"`
	Encoding          string `mapstructure:"encoding"`
	Level             string `mapstructure:"level"`
	Filename          string `mapstructure:"filename"`
	FileMaxSize       int    `mapstructure:"fileMaxSize"`
	FileMaxAge        int    `mapstructure:"fileMaxAge"`
	FileMaxBackups    int    `mapstructure:"fileMaxBackups"`
	FileIsCompress    bool   `mapstructure:"fileIsCompress"`
}

type topic struct {
	Name              string `mapstructure:"name"`
	Partition         int    `mapstructure:"partition"`
	ReplicationFactor int    `mapstructure:"replicationFactor"`
}

type KafkaConfig struct {
	Brokers []string         `mapstructure:"brokers"`
	Topics  map[string]topic `mapstructure:"topics"`
	GroupID string           `mapstructure:"groupID"`
}

type JaegerConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"serviceName"`
	LogSpans    bool   `mapstructure:"logSpans"`
}

func GetConf() (*Config, error) {
	re := regexp.MustCompile(`\$\{([^{}]+)\}`)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if re.Match([]byte(value)) {
			env := string(re.ReplaceAll([]byte(value), []byte("$1")))
			viper.Set(k, os.Getenv(env))
		}

	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	config.Server.ContextTimeout = config.Server.ContextTimeout * time.Second

	return &config, nil
}
