package api

import (
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/components/item"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/components/list"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "pong",
			})
	})

	// loads all html in templates dir
	r.LoadHTMLGlob("templates/*")

	r.GET("/home", func(c *gin.Context) {
		list.Render(c)
	})

	r.GET("/star/:id", func(c *gin.Context) {
		item.Render(c)
	})

	r.Static("/image", "./resources")
	r.Static("/styles", "./styles")

	// listen and serve on 127.0.0.1:8080
	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
