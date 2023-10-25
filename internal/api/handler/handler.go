package handler

import (
	"log"
	"net/http"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo api.Repo
}

func NewHandler(repo api.Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.GET("/ping", h.Ping)

	starRouter := r.Group("star") 
	{
		starRouter.GET("/", h.GetStarList)
		starRouter.GET("/:id", h.GetStar)
		starRouter.POST("/", h.CreateStar)
		starRouter.PUT("/:id/delete", h.DeleteStar)
		starRouter.PUT("/:id/update", h.UpdateStar)
		starRouter.PUT("/:id/event", h.PutIntoEvent)
	}

	eventRouter := r.Group("event")
	{
		eventRouter.GET("/", h.GetEventList)
		eventRouter.GET("/:id", h.GetEvent)
		eventRouter.PUT("/:id/update", h.UpdateEvent)
		eventRouter.PUT("/create", h.CreateEvent)
		eventRouter.PUT("/:id/form", h.FormEvent)
		eventRouter.PUT("/:id/complete", h.CompleteEvent)
		eventRouter.PUT("/:id/reject", h.RejectEvent)
		eventRouter.DELETE("/:id", h.DeleteEvent)
	}

	starEventRouter := r.Group("star-event")
	{
		starEventRouter.DELETE("/:id", h.RemoveFromEvent)
	}

	// listen and serve on 127.0.0.1:8080
	err := r.Run()
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
