package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorF struct {
	Message string `json:"error"`
}

type statusResponse struct {
	Status string `json:"Status"`
}

//	type jsonResponse struct {
//		Data []models.ProductVariant `json:"data"`
//	}
//
//	type jsonResponsePV struct {
//		Data *[]models.ProductVariant `json:"data"`
//	}
//
//	type jsonResponseCC struct {
//		Data *[]models.Condition `json:"conditions"`
//	}
//
//	type jsonResponseHistory struct {
//		Data *models.HistoryOutput `json:"data"`
//	}
type jsonResponseInterface struct {
	Data map[string]interface{} `json:"data"`
}
type jsonResponseInterface2 struct {
	Data interface{} `json:"data"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Errorf(message)
	c.AbortWithStatusJSON(statusCode, map[string]interface{}{
		"error": message,
	})
}
