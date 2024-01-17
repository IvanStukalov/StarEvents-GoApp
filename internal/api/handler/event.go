package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"errors"
	"StarEvent-GoApp/internal/models"

	"github.com/gin-gonic/gin"
)

const (
	Token      = "rjbrbjwf4908h8nfieh"
	ServiceUrl = "http://127.0.0.1:8081/scan/"
)

// GetEventList godoc
//
//	@Summary		Получить список событий
//	@Description	Возвращает список событий, отфильтрованных по заданным параметрам
//	@Tags			События
//	@Accept			json
//	@Produce		json
//	@Param			status			query		string			false	"Статус события"
//	@Param			start_formation	query		string			false	"Верхняя граница формирования события"
//	@Param			end_formation	query		string			false	"Нижняя граница формирования события"
//	@Success		200				{array}		models.Event	"Список событий"
//	@Failure		400				{string}	string			"Некорректный формат даты"
//	@Failure		404				{string}	string			"События не найдены"
//	@Router			/api/event/ [get]
func (h *Handler) GetEventList(c *gin.Context) {
	var startFormation time.Time
	var endFormation time.Time
	var err error

	status := c.Query("status")

	if c.Query("start_formation") != "" {
		startFormation, err = time.Parse(time.DateTime, c.Query("start_formation"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			log.Println(err)
			return
		}
	}

	if c.Query("end_formation") != "" {
		endFormation, err = time.Parse(time.DateTime, c.Query("end_formation"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			log.Println(err)
			return
		}
	}

	var creatorId = c.GetInt(userCtx)
	var isAdmin = c.GetBool(adminCtx)

	eventList, err := h.repo.GetEventList(status, startFormation, endFormation, creatorId, isAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, eventList)
}

// GetEvent godoc
//
//	@Summary		Получить событие по ID
//	@Description	Возвращает информацию о событии по его ID
//	@Tags			События
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"ID события"
//	@Success		200	{object}	map[string]interface{}	"Информация о событии"
//	@Failure		400	{string}	string			"Некорректный ID события"
//	@Failure		404	{string}	string			"Событие не найдено"
//	@Router			/api/event/{id} [get]
func (h *Handler) GetEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	event, starList, err := h.repo.GetEventByID(id, c.GetInt(userCtx), c.GetBool(adminCtx))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event, "star_list": starList})
}

// UpdateEvent godoc
//
//	@Summary		Обновить событие
//	@Description	Обновляет существующее событие по его ID
//	@Tags			События
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"ID события"
//	@Param			name	query		string	true	"Новое название события"
//	@Success		200		{string}	string	"Событие успешно обновлено"
//	@Failure		400		{string}	string	"Некорректный ID события или название"
//	@Failure		500		{string}	string	"Ошибка сервера"
//	@Router			/api/event/{id} [put]
func (h *Handler) UpdateEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	name := c.Query("name")
	err = h.repo.UpdateEvent(eventId, name)
}

// DeleteEvent godoc
//
//	@Summary		Удалить событие
//	@Description	Удаляет существующее событие
//	@Tags			События
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"Событие успешно удалено"
//	@Failure		400	{string}	string	"Ошибка удаления события"
//	@Failure		500	{string}	string	"Ошибка сервера"
//	@Router			/api/event/ [delete]
func (h *Handler) DeleteEvent(c *gin.Context) {
	creatorId := c.GetInt(userCtx)
	err := h.repo.DeleteEvent(creatorId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "событие успешно удалено"})
}

// FormEvent godoc
//
//	@Summary		Создать событие
//	@Description	Создает новое событие
//	@Tags			События
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"Событие успешно создано"
//	@Failure		400	{string}	string	"Ошибка создания события"
//	@Failure		500	{string}	string	"Ошибка сервера"
//	@Router			/api/event/form [put]
func (h *Handler) FormEvent(c *gin.Context) {
	err, eventId := h.repo.FormEvent(c.GetInt(userCtx))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.StartScanning(eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"mesage": err})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус изменен"})
}

// ChangeEventStatus godoc
//
//	@Summary		Изменить статус события
//	@Description	Изменяет статус существующего события
//	@Tags			События
//	@Accept			json
//	@Produce		json
//	@Param			status	query		string	true	"Новый статус события"
//	@Param			id		path		int		true	"ID события"
//	@Success		200		{string}	string	"Статус успешно изменен"
//	@Failure		400		{string}	string	"Некорректный статус или ID события"
//	@Failure		500		{string}	string	"Ошибка сервера"
//	@Router			/api/event/{id}/status [put]
func (h *Handler) ChangeEventStatus(c *gin.Context) {
	status := c.Query("status")
	if status != models.StatusAccepted && status != models.StatusCanceled {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Поменять статус можно только на 'Принято' и 'Отклонено'"})
		return
	}

	eventIdStr := c.Param("id")
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	err = h.repo.ChangeEventStatus(eventId, status, c.GetInt(userCtx))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Статус изменен"})
	return
}

func (h *Handler) StartScanning(eventId int) error {
	var event models.EventAsync
	event.ID = eventId
	body, _ := json.Marshal(event)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", ServiceUrl, bytes.NewBuffer(body))
	if err != nil {
		return errors.New("ошибка создания запроса")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return errors.New("ошибка запроса")
	}

	if resp.StatusCode == 200 {
		return nil
	}
	
	return errors.New("запрос не принят в обработку")
}

// FinishScanning godoc
// @Summary Завершает процесс сканирования
// @Description Принимает JSON запрос, проверяет токен и сохраняет процент сканирования. Возвращает сообщение об успешности или ошибке.
// @Tags			События
//
// @Accept json
// @Produce json
// @Param event body map[string]interface{} true "Событие для завершения сканирования"
// @Success 200 {object} map[string]string "Сообщение об успешности"
// @Failure 400 {object} map[string]string "Ошибка при привязке JSON"
// @Failure 403 {object} map[string]string "Неверный токен"
// @Router /api/event/finish-scanning [put]
func (h *Handler) FinishScanning(c *gin.Context) {
	var event models.EventAsync
	if err := c.BindJSON(&event); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	if event.Token != Token {
		c.AbortWithError(http.StatusForbidden, gin.H{"message": "неверный токен"})
		return
	}

	err := h.repo.SaveScannedPercent(event)
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "данные сохранены"})
}
