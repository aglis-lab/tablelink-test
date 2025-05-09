package app

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type appContext struct {
	config      *Configuration
	gorm        *gorm.DB
	RedisClient *redis.Client
	validator   *validator.Validate
	tracer      trace.Tracer
}

var appCtx appContext

func Init(ctx context.Context) error {
	// Init Config
	config, err := InitConfig(ctx)
	if err != nil {
		return err
	}

	// Init Postgresql Database
	configDb := config.Postgres
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable", configDb.Username, configDb.Password, configDb.Database)
	gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	// Set Gorm Tracing
	// if err := gorm.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
	// 	log.Panic(err)
	// }

	rdb, err := NewRedis()
	if err != nil {
		log.Panic(err)
	}

	appCtx = appContext{
		gorm:        gormDb,
		RedisClient: rdb,
		config:      config,
		tracer:      otel.Tracer(config.ServiceName, trace.WithInstrumentationVersion(config.ServiceVersion)),
		validator:   validator.New(),
	}

	return nil
}

func GormDB() *gorm.DB {
	return appCtx.gorm
}

func Config() *Configuration {
	return appCtx.config
}

func Tracer() trace.Tracer {
	return appCtx.tracer
}

func Validator() *validator.Validate {
	return appCtx.validator
}
