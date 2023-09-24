package list

import (
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/models"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/render"
	"github.com/gin-gonic/gin"
)

func Render(url string, c *gin.Context) {
	files := []string{
		"templates/list.gohtml",
	}

	data := models.GetData()

	render.RenderTmpl(url, files, data, c)
}
