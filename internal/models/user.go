package models

import (
	"errors"
)

type User struct {
	UserId           int    `gorm:"primaryKey" json:"userId"`
	Login            string `json:"login" binding:"required,max=64"`
	IsAdmin          bool   `json:"isAdmin" gorm:"column:is_moderator"`
	Password         string `json:"password,omitempty" binding:"required,min=8,max=64"`
}

type UserLogin struct {
	Login    string `json:"login" binding:"required,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type UserSignUp struct {
	Login    string `json:"login" binding:"required,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

var (
	ErrClientAlreadyExists = errors.New("клиент с таким логином уже существует")
	ErrUserNotFound        = errors.New("клиента с таким логином не существует")
)
