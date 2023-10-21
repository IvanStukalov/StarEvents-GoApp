package handler

import (
	"log"
	"net/http"
	"strconv"

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