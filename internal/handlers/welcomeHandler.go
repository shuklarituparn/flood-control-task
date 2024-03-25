package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WelcomeHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "welcome.html", gin.H{})
}
