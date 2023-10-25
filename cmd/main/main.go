package main

import (
	"time"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/handler"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/repository"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	minio "github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg/minio"
)

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

	h := handler.NewHandler(repo)
	h.StartServer()
}

func initConfig(vp *viper.Viper) error {
	vp.AddConfigPath("./config")
	vp.SetConfigName("config")

	return vp.ReadInConfig()
}
