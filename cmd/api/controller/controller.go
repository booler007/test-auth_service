package controller

import (
	"net/http"

	"authentication_medods/cmd/api/service"

	"github.com/gin-gonic/gin"
)

type Servicer interface {
	Authenticate(string, string) (*service.Tokens, error)
	Refresh()
}

type APIController struct {
	Service Servicer
}

func NewAPIController(s Servicer) *APIController {
	return &APIController{s}
}

func (c *APIController) SetupRouter(r *gin.Engine) {
	apiv1 := r.Group("/api/v1")

	apiv1.POST("/auth/signin/:uuid", c.SignIn)
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

}
