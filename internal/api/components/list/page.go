package list

import (
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/api/render"
	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context) {
	files := []string{
		"templates/list.tmpl",
	}

	data := item {
		title:
	}

	render.RenderTmpl("/home", files, data, c)
}
