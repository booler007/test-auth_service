package controller

import (
	"errors"
	"net/http"

	"authentication_medods/cmd/api/service"

	"github.com/gin-gonic/gin"
)

type Servicer interface {
	Authenticate(string, string) (*service.Tokens, error)
	RefreshTokens(string, string) (*service.Tokens, error)
}

type APIController struct {
	Service Servicer
}

type inputRefresh struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

var ErrWithBind = errors.New("error occurred while binding the request body")

func NewAPIController(s Servicer) *APIController {
	return &APIController{s}
}

func (c *APIController) SetupRouter(r *gin.Engine) {
	apiv1 := r.Group("/api/v1")

	apiv1.GET("/auth/signin/:uuid", c.SignIn)
	apiv1.POST("/auth/refresh", c.Refresh)
}

func (c *APIController) SignIn(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	tokens, err := c.Service.Authenticate(uuid, ctx.Request.RemoteAddr)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}

func (c *APIController) Refresh(ctx *gin.Context) {
	var input inputRefresh
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(ErrWithBind).SetType(gin.ErrorTypeBind)
		return
	}

	tokens, err := c.Service.RefreshTokens(input.RefreshToken, ctx.Request.RemoteAddr)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}
