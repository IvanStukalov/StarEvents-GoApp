package item

import (
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/models"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/render"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func Render(url string, c *gin.Context) {
	files := []string{
		"templates/item.gohtml",
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatalln("cannot convert string to int")
	}
	list := models.GetData()

	item := models.GetItemById(list, id)

	render.RenderTmpl(url, files, item, c)
}
