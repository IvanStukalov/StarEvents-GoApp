package handler

import (
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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

	r.GET("/home", h.GetStarList)
	r.GET("/star/:id", h.GetStarById)
	r.POST("/star/:id", h.DeleteStarById)

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

// get all stars
func (h *Handler) GetStarList(c *gin.Context) {
	data, err := h.repo.GetStarsByNameFilter(c.Query("starName"))
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

// get one star by id
func (h *Handler) GetStarById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		log.Println(err)
		return
	}

	item, err := h.repo.GetStarByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) DeleteStarById(c *gin.Context) {
	cardId := c.Param("id")
	id, err := strconv.Atoi(cardId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		return
	}

	err = h.repo.DeleteStarById(id)
	if err != nil { // если не получилось
		c.JSON(http.StatusBadRequest, nil)
		log.Printf("cant get star by id %v", err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
