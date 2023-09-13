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

	data := models.List{
		Items: []models.Item{
			{
				Name:        "Солнце",
				Description: "Наша родная звезда, которая светит нам и греет нас",
				Distance:    0,
				Magnitude:   -26.7,
				Image:       "sun.png",
			},
			{
				Name:        "Солнце",
				Description: "Наша родная звезда, которая светит нам и греет нас",
				Distance:    0,
				Magnitude:   -26.7,
				Image:       "sun.webp",
			},
		},
	}

	render.RenderTmpl(url, files, data, c)
}
