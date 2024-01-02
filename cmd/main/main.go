package main

import (
	"context"
	"os"
	"time"

	"StarEvent-GoApp/internal/api/handler"
	"StarEvent-GoApp/internal/api/repository"
	"StarEvent-GoApp/internal/pkg"
	minio "StarEvent-GoApp/internal/pkg/minio"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"StarEvent-GoApp/internal/pkg/auth"
	redis "StarEvent-GoApp/internal/pkg/redis"
)

//	@title			Star Events App
//	@version		1.0
//	@description	App for serving star events

//	@host		localhost:8080
//	@schemes	http
//	@BasePath	/
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
