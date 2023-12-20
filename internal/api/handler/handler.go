package handler

import (
	"log"
	"net/http"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/pkg/minio"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo  api.Repo
	minio minio.Client
}

func NewHandler(repo api.Repo, minio minio.Client) *Handler {
	return &Handler{repo: repo, minio: minio}
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

	// listen and serve on 127.0.0.1:8080
	err := r.Run("192.168.181.153:8080")
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
