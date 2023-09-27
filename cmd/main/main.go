package main

import (
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/handler"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/repository"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg"
	log "github.com/sirupsen/logrus"
)

func main() {
	dsn, err := pkg.GetConnectionString()
	if err != nil {
		log.Error(err)
	}
	log.Info(dsn)

	repo, err := repository.NewRepository(dsn)
	if err != nil {
		log.Error(err)
	}

	h := handler.NewHandler(repo)
	h.StartServer()
}
