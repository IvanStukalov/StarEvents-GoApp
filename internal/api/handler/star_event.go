package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// RemoveStarFromEvent godoc
//	@Summary		"Удалить звезду из события"
//	@Description	"Удаляет звезду из события по ее ID"
//	@Tags			Событие-Звезды
//	@Accept			json
//	@Produce		json
//	@Param			star-id	path		int				true	"ID звезды"
//	@Success		200		{object}	models.Event	"Событие после удаления звезды"
//	@Failure		400		{string}	string			"Некорректный ID звезды"
//	@Failure		500		{string}	string			"Ошибка сервера"
//	@Router			/api/star-event/{star-id} [delete]
func (h *Handler) RemoveStarFromEvent(c *gin.Context) {
	starIdStr := c.Param("star-id")
	starId, err := strconv.Atoi(starIdStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	creatorId := c.GetInt(userCtx)

	event, starList, err := h.repo.RemoveStarFromEvent(creatorId, starId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"event": event, "star_list": starList})
}
