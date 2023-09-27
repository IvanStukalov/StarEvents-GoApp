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

	var data models.List
	if c.Query("starName") != "" {
		data = models.GetItemByName(models.GetData(), c.Query("starName"))
	} else {
		data = models.GetData()
	}

	render.RenderTmpl(url, files, gin.H{"data": data, "QueryParam": c.Query("starName")}, c)
}
