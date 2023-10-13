package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
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

	starRouter.GET("/", h.GetStarList)
	starRouter.GET("/:id", h.GetStar)
	starRouter.PUT("/delete/:id", h.DeleteStar)
	starRouter.PUT("/update/:id", h.UpdateStar)
	starRouter.POST("/", h.CreateStar)


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
func (h *Handler) GetStar(c *gin.Context) {
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

// delete star
func (h *Handler) DeleteStar(c *gin.Context) {
	cardId := c.Param("id")
	id, err := strconv.Atoi(cardId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		return
	}

	err = h.repo.DeleteStarById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		log.Printf("cant get star by id %v", err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// update star
func (h *Handler) UpdateStar(c *gin.Context) {
	cardId := c.Param("id")
	id, err := strconv.Atoi(cardId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		return
	}

	newStar := models.Star{}
	newStar.ID = id
	newStar.Name = c.DefaultQuery("name", "")
	newStar.Description = c.DefaultQuery("description", "")

	distance, err := strconv.ParseFloat(c.DefaultQuery("distance", "-1"), 32)
		if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		return
	}
	newStar.Distance = float32(distance)

	age, err := strconv.ParseFloat(c.DefaultQuery("age", "-1"), 32)
		if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		return
	}
	newStar.Age = float32(age)

	magnitude, err := strconv.ParseFloat(c.DefaultQuery("magnitude", "100"), 32)
		if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		return
	}
	newStar.Magnitude = float32(magnitude)	
	newStar.Image = c.DefaultQuery("image", "")

	err = h.repo.UpdateStar(newStar)
	if err != nil { // если не получилось
		c.JSON(http.StatusBadRequest, nil)
		log.Printf("cant get star by id %v", err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// create star
func (h *Handler) CreateStar(c *gin.Context) {
	var data models.Star
  if err := c.BindJSON(&data); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert json"})
    return
	}
  
	err := h.repo.CreateStar(data)
	if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"err_msg": "something wrong"})
    return
	}

	c.JSON(http.StatusOK, nil)
}
