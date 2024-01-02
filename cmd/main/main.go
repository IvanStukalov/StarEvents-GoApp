package main

import (
	"context"
	"os"
	"time"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/handler"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/repository"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg"
	minio "github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg/minio"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
  "github.com/IvanStukalov/Term5-WebAppDevelopment/docs"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg/auth"
	redis "github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg/redis"
)

// @title ThreatMonitoringApp
// @version 1.0
// @description App for serving threats monitoring requests

// @host localhost:8080
// @schemes http
// @BasePath /
func main() {
	dsn, err := pkg.GetConnectionString()
	if err != nil {
		log.Error(err)
	}
	log.Info(dsn)

	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	vp := viper.New()
	if err := initConfig(vp); err != nil {
		logger.Fatalf("error initializing configs: %s", err.Error())
	}

	minioConfig := minio.InitConfig(vp)

	minioClient, err := minio.NewMinioClient(context.Background(), minioConfig, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	repo, err := repository.NewRepository(dsn)
	if err != nil {
		log.Error(err)
	}

	redisConfig := redis.InitRedisConfig(vp, logger)

	redisClient, err := redis.NewRedisClient(context.Background(), redisConfig, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	tokenManager, err := auth.NewManager(os.Getenv("TOKEN_SECRET"))
	if err != nil {
		logger.Fatalln(err)
	}

	h := handler.NewHandler(repo, minioClient, tokenManager, redisClient)
	h.StartServer()
}

func initConfig(vp *viper.Viper) error {
	vp.AddConfigPath("./config")
	vp.SetConfigName("config")

	return vp.ReadInConfig()
}
