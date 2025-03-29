package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIHandler struct {
	Db *gorm.DB
}

func NewHandler(db *gorm.DB) *APIHandler {
	return &APIHandler{Db: db}
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
