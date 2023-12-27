package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
	"github.com/gin-gonic/gin"
)

const (
	Token = "rjbrbjwf4908h8nfieh"
	ServiceUrl  = "http://127.0.0.1:8081/scan/"
)

// get list of events
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

	eventList, err := h.repo.GetEventList(status, startFormation, endFormation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, eventList)
}

// get one event by id
func (h *Handler) GetEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	event, starList, err := h.repo.GetEventByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event, "star_list": starList})
}

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

// put star into event
func (h *Handler) PutIntoEvent(c *gin.Context) {
	var eventMsg models.EventMsg

	err := c.BindJSON(&eventMsg)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	eventMsg.CreatorID = h.repo.GetCreatorId()
	log.Println(eventMsg)

	err = h.repo.PutIntoEvent(eventMsg)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// deleting event with status
func (h *Handler) DeleteEvent(c *gin.Context) {
	creatorId := h.repo.GetCreatorId()
	err := h.repo.DeleteEvent(creatorId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "событие успешно удалено"})
}

func (h *Handler) FormEvent(c *gin.Context) {
	err := h.repo.FormEvent(h.repo.GetCreatorId())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус изменен"})
}

func (h *Handler) ChangeEventStatus(c *gin.Context) {
	status := c.Query("status")
	if status != models.StatusAccepted && status != models.StatusCanceled && status != models.StatusClosed {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Поменять статус можно только на 'accepted, 'closed' и 'canceled'"})
		return
	}

	eventIdStr := c.Param("id")
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	err = h.repo.ChangeEventStatus(eventId, status)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Статус изменен"})
	return
}

func (h *Handler) StartScanning(c *gin.Context) {
	var event models.EventAsync
	if err := c.BindJSON(&event); err != nil {
		c.AbortWithError(http.StatusBadRequest, gin.H{"message": "Неверный формат"})
		return
	}


	event.Token = Token
	fmt.Println("got event:", event)

	body, _ := json.Marshal(event)

	fmt.Println("Body formed:", body)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", ServiceUrl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	if resp.StatusCode == 200 {
		c.JSON(http.StatusOK, gin.H{"message": "заявка принята в обработку"})
		return
	}
	c.AbortWithError(http.StatusInternalServerError, gin.H{"message": "заявка не принята в обработку"})
}

func (h *Handler) FinishScanning(c *gin.Context) {
	var event models.EventAsync
	if err := c.BindJSON(&event); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	fmt.Println("Finish scanning:", event)

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
