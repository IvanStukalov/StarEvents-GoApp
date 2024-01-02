package handler

import (
	"errors"
	"fmt"
	"StarEvent-GoApp/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

// SignUp godoc
//	@Summary		Регистрация нового пользователя
//	@Description	Регистрирует нового пользователя с заданными параметрами
//	@Tags			Пользователи
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.UserSignUp	true	"Новый пользователь"
//	@Success		201		{string}	string				"Пользователь успешно зарегистрирован"
//	@Failure		400		{string}	string				"Неверный формат данных о новом пользователе"
//	@Failure		500		{string}	string				"Нельзя создать пользователя с таким логином"
//	@Router			/api/signUp [post]
func (h *Handler) SignUp(c *gin.Context) {
	var newClient models.UserSignUp
	var err error

	if err = c.BindJSON(&newClient); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных о новом пользователе"})
		return
	}

	if newClient.Password, err = h.hasher.Hash(newClient.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат пароля"})
		return
	}

	if err = h.repo.SignUp(c.Request.Context(), models.User{
		Login:    newClient.Login,
		Password: newClient.Password,
	}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "нельзя создать пользователя с таким логином"})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "пользователь успешно создан"})
}

// CheckAuth godoc
//	@Summary		Проверка аутентификации
//	@Description	Проверяет аутентификацию текущего пользователя и возвращает его информацию
//	@Tags			Аутентификация
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.User	"Информация о пользователе"
//	@Failure		500	{string}	string		"Ошибка сервера"
//	@Security		ApiKeyAuth
//	@Router			/api/check-auth [get]
func (h *Handler) CheckAuth(c *gin.Context) {
	var userInfo = models.User{UserId: c.GetInt(userCtx)}

	userInfo, err := h.repo.GetUserInfo(c.Request.Context(), userInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "неверный формат данных")
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

// SignIn godoc
//	@Summary		Авторизация пользователя
//	@Description	Авторизует пользователя и возвращает JWT токен
//	@Tags			Аутентификация
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.UserLogin	true	"Данные для входа"
//	@Success		200		{string}	string				"Пользователь успешно авторизован"
//	@Failure		400		{string}	string				"Неверный формат данных"
//	@Failure		401		{string}	string				"Неверные учетные данные"
//	@Failure		500		{string}	string				"Ошибка сервера"
//	@Security		ApiKeyAuth
//	@Router			/api/signIn [post]
func (h *Handler) SignIn(c *gin.Context) {
	var clientInfo models.UserLogin
	var err error

	if err = c.BindJSON(&clientInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "неверный формат данных")
		return
	}

	if clientInfo.Password, err = h.hasher.Hash(clientInfo.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "неверный формат пароля"})
		return
	}

	user, err := h.repo.GetByCredentials(c.Request.Context(), models.User{Password: clientInfo.Password, Login: clientInfo.Login})
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка авторизации"})
		return
	}

	token, err := h.tokenManager.NewJWT(user.UserId, user.IsAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ошибка при формировании токена"})
		return
	}

	c.SetCookie("AccessToken", "Bearer "+token, 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "клиент успешно авторизован"})
}

// Logout godoc
//	@Summary		Выход из системы
//	@Description	Отменяет аутентификацию пользователя и удаляет JWT токен
//	@Tags			Аутентификация
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"Пользователь успешно вышел из системы"
//	@Failure		400	{string}	string	"Неверный формат данных"
//	@Failure		500	{string}	string	"Ошибка сервера"
//	@Security		ApiKeyAuth
//	@Router			/api/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	jwtStr, err := c.Cookie("AccessToken")
	if !strings.HasPrefix(jwtStr, jwtPrefix) || err != nil { // если нет префикса то нас дурят!
		c.AbortWithStatus(http.StatusBadRequest) // отдаем что нет доступа
		return
	}

	// отрезаем префикс
	jwtStr = jwtStr[len(jwtPrefix):]

	_, _, err = h.tokenManager.Parse(jwtStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	// сохраняем в блеклист редиса
	err = h.redis.WriteJWTToBlacklist(c.Request.Context(), jwtStr, time.Hour)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
