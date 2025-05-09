package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type (
	MySQL struct {
		ConnURI            string        `mapstructure:"MYSQL_CONN_URI" validate:"required"`
		Database           string        `mapstructure:"MYSQL_DATABASE" validate:"required"`
		Username           string        `mapstructure:"MYSQL_USERNAME" validate:"required"`
		Password           string        `mapstructure:"MYSQL_PASSWORD"`
		MaxPoolSize        int           `mapstructure:"MYSQL_MAX_POOL_SZE"` //Optional, default to 0 (zero value of int)
		MaxIdleConnections int           `mapstructure:"MYSQL_MAX_IDLE_CONNECTIONS"`
		MaxIdleTime        time.Duration `mapstructure:"MYSQL_MAX_IDLE_TIME"` //Optional, default to '0s' (zero value of time.Duration)
		MaxLifeTime        time.Duration `mapstructure:"MYSQL_MAX_IDLE_TIME"` //Optional, default to '0s' (zero value of time.Duration)
	}

	Postgres struct {
		Database string `mapstructure:"POSTGRESQL_DATABASE" validate:"required"`
		Username string `mapstructure:"POSTGRESQL_USERNAME" validate:"required"`
		Password string `mapstructure:"POSTGRESQL_PASSWORD" validate:"required"`
	}

	Configuration struct {
		ServiceName    string `mapstructure:"SERVICE_NAME" validate:"required"`
		ServiceVersion string `mapstructure:"SERVICE_VERSION" validate:"required"`
		// MySQL          MySQL  `mapstructure:",squash"`
		Postgres Postgres `mapstructure:",squash"`

		OltpGRPCProvider string `mapstructure:"OLTP_GRPC_PROVIDER" validate:"required"`
		Environment      string `mapstructure:"ENV" validate:"required,oneof=development staging production"`
		BindAddress      int    `mapstructure:"BIND_ADDRESS" validate:"required"`
		GrpcPort         int    `mapstructure:"GRPC_PORT" validate:"required"`
	}
)

func InitConfig(ctx context.Context) (*Configuration, error) {
	var cfg Configuration

	viper.SetConfigType("env")
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	_, err := os.Stat(envFile)
	if !os.IsNotExist(err) {
		viper.SetConfigFile(envFile)

		if err := viper.ReadInConfig(); err != nil {
			log.Printf("failed to read config:%v", err)
			return nil, err
		}
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("failed to bind config:%v", err)
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Printf("invalid config:%v", err)
		}
		log.Printf("failed to load config")
		return nil, err
	}

	log.Printf("Config loaded: %+v", cfg)
	return &cfg, nil
}
