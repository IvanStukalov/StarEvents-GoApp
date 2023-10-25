package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
	"github.com/gin-gonic/gin"
)

// get list of events
func (h *Handler) GetEventList(c *gin.Context) {
	eventList, err := h.repo.GetEventList()
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, eventList)
}

// get one event by id
func (h *Handler) GetEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		log.Println(err)
		return
	}

	eventDetails, err := h.repo.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, eventDetails)
}

// update event by id
func (h *Handler) UpdateEvent(c *gin.Context) {
	var event models.Event
	var err error

	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert json"})
		return
	}

	event.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		log.Println(err)
		return
	}

	err = h.repo.UpdateEvent(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err_msg": "something wrong"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// create new event
func (h *Handler) CreateEvent(c *gin.Context) {
	event := models.Event{}

	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert json"})
		return
	}

	creator, err := h.repo.GetCreatorId()
	fmt.Println(creator)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err_msg": "cannot find creator"})
		return
	}

	moderator, err := h.repo.GetModeratorId()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err_msg": "cannot find moderator"})
		return
	}
	fmt.Println(moderator)

	event.CreatorID = creator
	event.ModeratorID = moderator
	event.Status = "pending"
	event.CreationDate = time.Now()

	newEvent, err := h.repo.CreateEvent(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err_msg": "cannot create event"})
		return
	}

	c.JSON(http.StatusOK, newEvent)
}

// form created event
func (h *Handler) FormEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		log.Println(err)
		return
	}

	formedEvent, err := h.repo.FormEvent(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err_msg": "cannot form event"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, formedEvent)
}

// complete event
func (h *Handler) CompleteEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		log.Println(err)
		return
	}

	completedEvent, err := h.repo.CompleteEvent(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err_msg": "cannot complete event"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, completedEvent)
}

// reject event
func (h *Handler) RejectEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		log.Println(err)
		return
	}

	rejectedEvent, err := h.repo.RejectEvent(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err_msg": "cannot reject event"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, rejectedEvent)
}

// delete event
func (h *Handler) DeleteEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err_msg": "cannot convert id to int"})
		log.Println(err)
		return
	}

	err = h.repo.DeleteEvent(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err_msg": "cannot delete event"})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
