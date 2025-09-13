package utils

import (
	"net/http"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func NewNotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func NewConflict(msg string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: msg}
}

func NewInternal(msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg}
}
