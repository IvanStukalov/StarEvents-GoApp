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
				Name:        "Проксима Центавра",
				Description: "Звезда, красный карлик, относящаяся к звёздной системе Альфа Центавра, ближайшая к Солнцу звезда",
				Distance:    4.2,
				Magnitude:   11.1,
				Image:       "Proxima_Centauri.jpg",
			},
			{
				Name:        "Звезда Барнарда",
				Description: "Одиночная звезда в созвездии Змееносца",
				Distance:    5.96,
				Magnitude:   9.57,
				Image:       "Barnard.jpeg",
			},
			{
				Name:        "Сириус",
				Description: "Ярчайшая звезда ночного неба",
				Distance:    8.6,
				Magnitude:   -1.46,
				Image:       "Sirius.jpg",
			},
			{
				Name:        "Лейтен 726-8",
				Description: "Двойная звезда в созвездии Кита",
				Distance:    8.73,
				Magnitude:   12.5,
				Image:       "Leiten.jpg",
			},
		},
	}

	render.RenderTmpl(url, files, data, c)
}
