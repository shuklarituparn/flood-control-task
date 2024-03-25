package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PersikHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"Привет": "Я персик",
	})
}
