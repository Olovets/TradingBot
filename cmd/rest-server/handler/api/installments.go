package api

import (
	"github.com/Olovets/TradingBot/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *APIHandler) Candles(c *gin.Context) {
	var r struct {
		Period    string `json:"period"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
	}
	var e error

	if err := c.BindJSON(&r); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	candles, _ := models.GetAllCandles(h.Db)

	res, e := models.GenerateCandles(candles, r.Period, r.StartTime, r.EndTime)

	if e != nil {
		newErrorResponse(c, http.StatusInternalServerError, e.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}
