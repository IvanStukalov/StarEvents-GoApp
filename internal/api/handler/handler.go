package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/swaggo/files"
	_ "StarEvent-GoApp/docs"
	"github.com/swaggo/gin-swagger"
	"StarEvent-GoApp/internal/api"
	"StarEvent-GoApp/internal/models"
	"StarEvent-GoApp/internal/pkg/auth"
	"StarEvent-GoApp/internal/pkg/hash"
	minio "StarEvent-GoApp/internal/pkg/minio"
	"StarEvent-GoApp/internal/pkg/redis"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	minio minio.Client
	redis redis.Client
	repo  api.Repo

	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
}

func NewHandler(repo api.Repo, minioClient minio.Client, tokenManager auth.TokenManager, redisClient redis.Client) *Handler {
	return &Handler{
		repo:         repo,
		minio:        minioClient,
		redis:        redisClient,
		hasher:       hash.NewSHA256Hasher(os.Getenv("SALT")),
		tokenManager: tokenManager,
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "localhost")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (h *Handler) StartServer() {
	log.Println("Server start up")

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	api := r.Group("api")
	api.GET("/ping", h.Ping)

	api.POST("/signIn", h.SignIn)
	api.POST("/signUp", h.SignUp)
	api.POST("/logout", h.Logout)
	api.GET("/check-auth", h.WithAuthCheck([]models.Role{models.Client, models.Admin}, false), h.CheckAuth)

	starRouter := api.Group("star")
	{
		starRouter.GET("/", h.WithAuthCheck([]models.Role{models.Admin, models.Client}, true), h.GetStarList)
		starRouter.GET("/:id", h.GetStar)
		starRouter.POST("/", h.WithAuthCheck([]models.Role{models.Admin}, false), h.CreateStar)
		starRouter.DELETE("/:id", h.WithAuthCheck([]models.Role{models.Admin}, false), h.DeleteStar)
		starRouter.PUT("/:id/update", h.WithAuthCheck([]models.Role{models.Admin}, false), h.UpdateStar)
		starRouter.POST("/event", h.WithAuthCheck([]models.Role{models.Client}, false), h.PutIntoEvent)
	}

	eventRouter := api.Group("event")
	{
		eventRouter.GET("/", h.WithAuthCheck([]models.Role{models.Admin, models.Client}, false), h.GetEventList)
		eventRouter.GET("/:id", h.WithAuthCheck([]models.Role{models.Admin, models.Client}, false), h.GetEvent)
		eventRouter.PUT("/:id", h.WithAuthCheck([]models.Role{models.Client}, false), h.UpdateEvent)
		eventRouter.DELETE("/", h.WithAuthCheck([]models.Role{models.Client}, false), h.DeleteEvent)
		eventRouter.PUT("/form", h.WithAuthCheck([]models.Role{models.Client}, false), h.FormEvent)
		eventRouter.PUT("/:id/status", h.WithAuthCheck([]models.Role{models.Admin}, false), h.ChangeEventStatus)

		// async
		eventRouter.PUT("/start-scanning", h.WithAuthCheck([]models.Role{models.Admin}, false), h.StartScanning)
		eventRouter.PUT("/finish-scanning", h.FinishScanning)
	}

	starEventRouter := api.Group("star-event")
	{
		starEventRouter.DELETE("/:star-id", h.WithAuthCheck([]models.Role{models.Client}, false), h.RemoveStarFromEvent)
	}

	err := r.Run(fmt.Sprintf("%s:8080", os.Getenv("HOST")))
	if err != nil {
		log.Fatalln(err)
	}
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "pong",
		})
}
