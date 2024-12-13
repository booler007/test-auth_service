package controller

import (
	"errors"
	"log"
	"net/http"

	"authentication_medods/cmd/api/service"
	"authentication_medods/cmd/api/storage"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if err := ctx.Errors.Last(); err != nil {
			switch {
			case err.Type == gin.ErrorTypeBind:
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case errors.Is(err, service.ErrInvalidRefreshToken):
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case errors.Is(err, storage.ErrUserNotFound):
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			default:
				log.Print(err.Error())
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}
	}
}
