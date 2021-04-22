package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type successfulResponse struct {
	Data interface{} `json:"data"`
}

func formSuccessResponse(status int, c *gin.Context, data interface{}) {
	c.JSON(status, &successfulResponse{Data: data})
}

func formOkResponse(c *gin.Context, data interface{}) {
	formSuccessResponse(http.StatusOK, c, data)
}

func formCreatedResponse(c *gin.Context, data interface{}) {
	formSuccessResponse(http.StatusCreated, c, data)
}
