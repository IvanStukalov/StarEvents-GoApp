package models

import "strings"

func GetData() List {
	return List{
		Items: []Item{
			{
				ID:          1,
				Name:        "Солнце",
				Description: "Наша родная звезда, которая светит нам и греет нас",
				Distance:    0,
				Magnitude:   -26.7,
				Image:       "sun.png",
			},
			{
				ID:          2,
				Name:        "Проксима Центавра",
				Description: "Звезда, красный карлик, относящаяся к звёздной системе Альфа Центавра, ближайшая к Солнцу звезда",
				Distance:    4.2,
				Magnitude:   11.1,
				Image:       "Proxima_Centauri.jpg",
			},
			{
				ID:          3,
				Name:        "Звезда Барнарда",
				Description: "Одиночная звезда в созвездии Змееносца",
				Distance:    5.96,
				Magnitude:   9.57,
				Image:       "Barnard.jpeg",
			},
			{
				ID:          4,
				Name:        "Сириус",
				Description: "Ярчайшая звезда ночного неба",
				Distance:    8.6,
				Magnitude:   -1.46,
				Image:       "Sirius.jpg",
			},
			{
				ID:          5,
				Name:        "Лейтен 726-8",
				Description: "Двойная звезда в созвездии Кита",
				Distance:    8.73,
				Magnitude:   12.5,
				Image:       "Leiten.jpg",
			},
		},
	}
}

func GetItemById(list List, id int) Item {
	for i := 0; i < len(list.Items); i++ {
		if id == list.Items[i].ID {
			return list.Items[i]
		}
	}
	return Item{}
}

func GetItemByName(list List, name string) List {
	res := List{}
	for i := 0; i < len(list.Items); i++ {
		if strings.Contains(strings.ToLower(list.Items[i].Name), strings.ToLower(name)) {
			res.Items = append(res.Items, list.Items[i])
		}
	}
	return res
}
