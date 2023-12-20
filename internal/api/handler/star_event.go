package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) RemoveStarFromEvent(c *gin.Context) {
	starIdStr := c.Param("star-id")
	starId, err := strconv.Atoi(starIdStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	creatorId := h.repo.GetCreatorId()

	event, starList, err := h.repo.RemoveStarFromEvent(creatorId, starId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"event": event, "star_list": starList})
}
