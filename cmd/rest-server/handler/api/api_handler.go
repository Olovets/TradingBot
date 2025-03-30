package api

import (
	"github.com/gin-gonic/gin"
)

type APIHandler struct {
}

func NewHandler() *APIHandler {
	return &APIHandler{}
}

func (h *APIHandler) InitAPIRoutes(apiRouter *gin.RouterGroup) {

	api_ := apiRouter.Group("/api")
	{
		candles := api_.Group("/candles")
		{
			candles.POST("", h.Candles)
		}
	}

}
