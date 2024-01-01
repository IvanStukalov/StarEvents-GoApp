package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/repository"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg/auth"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg/hash"
	minio "github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg/minio"
	redis "github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg/redis"
)

type Handler struct {
	minio minio.Client
	redis redis.Client
	repo  api.Repo

	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
}

func NewHandler(repo api.Repo, minio minio.Client, tokenManager, redisClient) *Handler {
	return &Handler{
		repo:         repo,
		minio:        minioClient,
		redis:        redisClient,
		hasher:       hash.NewSHA256Hasher(os.Getenv("SALT")),
		tokenManager: tokenManager,
	}
}

func (h *Handler) StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	api := r.Group("api")

	api.GET("/ping", h.Ping)

	starRouter := api.Group("star")
	{
		starRouter.GET("/", h.GetStarList)
		starRouter.GET("/:id", h.GetStar)
		starRouter.POST("/", h.CreateStar)
		starRouter.DELETE("/:id", h.DeleteStar)
		starRouter.PUT("/:id/update", h.UpdateStar)
	}

	eventRouter := api.Group("event")
	{
		eventRouter.GET("/", h.GetEventList)
		eventRouter.GET("/:id", h.GetEvent)
		eventRouter.PUT("/:id", h.UpdateEvent)
		eventRouter.POST("/star", h.PutIntoEvent)
		eventRouter.DELETE("/", h.DeleteEvent)
		eventRouter.PUT("/form", h.FormEvent)
		eventRouter.PUT("/:id/status", h.ChangeEventStatus)
	}

	starEventRouter := api.Group("star-event")
	{
		starEventRouter.DELETE("/:star-id", h.RemoveStarFromEvent)
	}

	err := r.Run(fmt.Sprintf("%s:8080", os.Getenv("HOST")))
	if err != nil {
		log.Fatalln(err)
	}
}

// ping
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "pong",
		})
}
