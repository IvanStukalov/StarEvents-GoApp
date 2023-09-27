package handler

import (
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/render"
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

	// loads all html in templates dir
	r.LoadHTMLGlob("templates/*")

	r.GET("/home", h.GetStarList)
	r.GET("/star/:id", h.GetStarById)

	r.Static("/image", "./resources")
	r.Static("/styles", "./styles")

	// listen and serve on 127.0.0.1:8080
	err := r.Run()
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

func (h *Handler) GetStarList(c *gin.Context) {
	files := []string{
		"templates/list.gohtml",
	}

	//var data []models.Star
	//if c.Query("name") != "" {
	//	data = models.GetItemByName(models.GetData(), c.Query("name"))
	//} else {
	data, err := h.repo.GetStars()
	if err != nil {
		log.Println(err)
	}

	render.RenderTmpl(files, data, c)
}

func (h *Handler) GetStarById(c *gin.Context) {
	files := []string{
		"templates/item.gohtml",
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
	}

	item, err := h.repo.GetStarByID(id)
	if err != nil {
		log.Println(err)
	}

	render.RenderTmpl(files, item, c)
}
