package handler

import (
	"log"
	"net/http"
	"strconv"

	"StarEvent-GoApp/internal/models"
	"StarEvent-GoApp/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetStarList godoc
//	@Summary		Получить список звезд
//	@Description	Возвращает список звезд, отфильтрованных по заданным параметрам
//	@Tags			Звезды
//	@Accept			json
//	@Produce		json
//	@Param			name		query		string					false	"Имя"
//	@Param			dist_top	query		float64					false	"Расстояние до верхней границы"
//	@Param			dist_bot	query		float64					false	"Расстояние до нижней границы"
//	@Param			age_top		query		int						false	"Возраст до верхней границы"
//	@Param			age_bot		query		int						false	"Возраст до нижней границы"
//	@Param			mag_top		query		float64					false	"Магнитуда до верхней границы"
//	@Param			mag_bot		query		float64					false	"Магнитуда до нижней границы"
//	@Success		200			{object}	map[string]interface{}	"Успешный ответ"
//	@Failure		404			{object}	map[string]interface{}	"Ошибка при получении списка звезд или черновика"
//	@Router			/api/star [get]
func (h *Handler) GetStarList(c *gin.Context) {
	starList, err := h.repo.GetFilteredStars(c.Query("name"), 
																					 c.Query("dist_top"), 
																					 c.Query("dist_bot"), 
																					 c.Query("age_top"), 
																					 c.Query("age_bot"), 
																					 c.Query("mag_top"), 
																					 c.Query("mag_bot"))
																					 
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, nil)
		log.Println(err)
		return
	}

	draftId, err := h.repo.GetDraft(c.GetInt(userCtx))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, nil)
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"stars": starList, "draft_id": draftId})
}

// GetStar godoc
//	@Summary		Получить звезду по ID
//	@Description	Возвращает информацию о звезде по ее ID
//	@Tags			Звезды
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int			true	"ID звезды"
//	@Success		200	{object}	models.Star	"Успешный ответ"
//	@Failure		400	{string}	string		"Некорректный ID звезды"
//	@Failure		404	{string}	string		"Звезда не найдена"
//	@Router			/api/star/{id} [get]
func (h *Handler) GetStar(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	star, err := h.repo.GetStarByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, star)
}

// CreateStar godoc
//	@Summary		Создать звезду
//	@Description	Создает новую звезду с заданными параметрами
//	@Tags			Звезды
//	@Accept			mpfd
//	@Produce		json
//	@Param			name		formData	string	true	"Название звезды"
//	@Param			description	formData	string	false	"Описание звезды"
//	@Param			distance	formData	float32	false	"Расстояние до звезды"
//	@Param			age			formData	float32	false	"Возраст звезды"
//	@Param			magnitude	formData	float32	false	"Магнитуда звезды"
//	@Param			image		formData	file	true	"Изображение звезды"
//	@Success		200			{string}	string	"Успешное создание звезды"
//	@Failure		400			{string}	string	"Некорректный ввод данных"
//	@Failure		500			{string}	string	"Ошибка сервера"
//	@Router			/api/star [post]
func (h *Handler) CreateStar(c *gin.Context) {
	var star models.Star

	star.Name = c.Request.FormValue("name")
	star.Description = c.Request.FormValue("description")

	distanceValue := c.Request.FormValue("distance")
	if distanceValue != "" {
		distance, err := strconv.ParseFloat(distanceValue, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			log.Println(err)
			return
		}
		star.Distance = float32(distance)
	}

	ageValue := c.Request.FormValue("age")
	if ageValue != "" {
		age, err := strconv.ParseFloat(ageValue, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			log.Println(err)
			return
		}
		star.Age = float32(age)
	}

	magnitudeValue := c.Request.FormValue("magnitude")
	if magnitudeValue != "" {
		magnitude, err := strconv.ParseFloat(magnitudeValue, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			log.Println(err)
			return
		}
		star.Magnitude = float32(magnitude)
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	if star.Image, err = h.minio.SaveImage(c.Request.Context(), file, header); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	err = h.repo.CreateStar(star)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// UpdateStar godoc
//	@Summary		Обновить звезду
//	@Description	Обновляет существующую звезду с заданными параметрами
//	@Tags			Звезды
//	@Accept			mpfd
//	@Produce		json
//	@Param			id			path		int		true	"ID звезды"
//	@Param			name		formData	string	false	"Новое название звезды"
//	@Param			description	formData	string	false	"Новое описание звезды"
//	@Param			distance	formData	float32	false	"Новое расстояние до звезды"
//	@Param			age			formData	float32	false	"Новый возраст звезды"
//	@Param			magnitude	formData	float32	false	"Новая магнитуда звезды"
//	@Param			image		formData	file	false	"Новое изображение звезды"
//	@Success		200			{string}	string	"Успешное обновление звезды"
//	@Failure		400			{string}	string	"Некорректный ввод данных"
//	@Failure		500			{string}	string	"Ошибка сервера"
//	@Router			/api/star/{id}/update [put]
func (h *Handler) UpdateStar(c *gin.Context) {
	var updatedStar models.Star
	var isUpdA bool
	var isUpdD bool
	var isUpdM bool

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	updatedStar.ID = id
	updatedStar.Name = c.Request.FormValue("name")
	updatedStar.Description = c.Request.FormValue("description")

	distanceValue := c.Request.FormValue("distance")
	if distanceValue != "" {
		isUpdD = true
		distance, err := strconv.ParseFloat(distanceValue, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			log.Println(err)
			return
		}
		updatedStar.Distance = float32(distance)
	}

	ageValue := c.Request.FormValue("age")
	if ageValue != "" {
		isUpdA = true
		age, err := strconv.ParseFloat(ageValue, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			log.Println(err)
			return
		}
		updatedStar.Age = float32(age)
	}

	magnitudeValue := c.Request.FormValue("magnitude")
	if magnitudeValue != "" {
		isUpdM = true
		magnitude, err := strconv.ParseFloat(magnitudeValue, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			log.Println(err)
			return
		}
		updatedStar.Magnitude = float32(magnitude)
	}

	file, header, err := c.Request.FormFile("image")

	if header != nil && header.Size != 0 {
		if updatedStar.Image, err = h.minio.SaveImage(c.Request.Context(), file, header); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		// delete old image from db
		url, err := h.repo.GetStarImageById(updatedStar.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			log.Println(err)
			return
		}
		// delete image from minio
		err = h.minio.DeleteImage(c.Request.Context(), utils.ExtractObjectNameFromUrl(url))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			log.Println(err)
			return
		}
	}

	err = h.repo.UpdateStar(updatedStar, isUpdA, isUpdD, isUpdM)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// DeleteStar godoc
//	@Summary		Удалить звезду
//	@Description	Удаляет существующую звезду по ее ID
//	@Tags			Звезды
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"ID звезды"
//	@Success		200	{string}	string	"Успешное удаление звезды"
//	@Failure		400	{string}	string	"Некорректный ID звезды"
//	@Failure		500	{string}	string	"Ошибка сервера"
//	@Router			/api/star/{id} [delete]
func (h *Handler) DeleteStar(c *gin.Context) {
	cardId := c.Param("id")
	id, err := strconv.Atoi(cardId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	err = h.repo.DeleteStarByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	// delete old image from db
	url, err := h.repo.GetStarImageById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	// delete image from minio
	err = h.minio.DeleteImage(c.Request.Context(), utils.ExtractObjectNameFromUrl(url))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// PutIntoEvent godoc
//	@Summary		Добавить сообщение в событие
//	@Description	Добавляет сообщение в событие по его ID
//	@Tags			События
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"ID события"
//	@Param			message	body		models.EventMsg	true	"Сообщение для добавления в событие"
//	@Success		200		{string}	string			"Сообщение успешно добавлено в событие"
//	@Failure		400		{string}	string			"Некорректный ID события или сообщение"
//	@Failure		500		{string}	string			"Ошибка сервера"
//	@Router			/api/star/event [post]
func (h *Handler) PutIntoEvent(c *gin.Context) {
	var eventMsg models.EventMsg

	err := c.BindJSON(&eventMsg)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	eventMsg.CreatorID = c.GetInt(userCtx)
	log.Println(eventMsg)

	err = h.repo.PutIntoEvent(eventMsg)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
