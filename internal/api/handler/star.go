package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/models"
	"github.com/IvanStukalov/Term5-WebAppDevelopment/internal/utils"

	"github.com/gin-gonic/gin"
)

// get stars with filter
func (h *Handler) GetStarList(c *gin.Context) {
	starList, err := h.repo.GetFilteredStars(c.Query("name"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, nil)
		log.Println(err)
		return
	}

	draftId, err := h.repo.GetDraft(h.repo.GetCreatorId())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, nil)
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"stars": starList, "draft_id": draftId})
}

// get one star by id
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

// create star
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

// update star
func (h *Handler) UpdateStar(c *gin.Context) {
	var updatedStar models.Star

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

	err = h.repo.UpdateStar(updatedStar)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// delete star
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
