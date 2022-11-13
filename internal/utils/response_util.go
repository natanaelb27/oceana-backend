package utils

import (
	"net/http"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

func newStandardError(message string, data any) gin.H {
	return gin.H{
		"message": message,
		"data":    data,
	}
}

func formatError(err error) (standardErr gin.H) {
	if err != nil {
		newErr := TranslateError(err)
		var message string
		if len(newErr) > 0 {
			message = "validation"
		} else {
			message = err.Error()
		}
		standardErr = newStandardError(message, newErr)
	}
	return
}

func response(c *gin.Context, statusCode int, data any, metadata any, err error, message string) dtos.Response {
	standardErr := formatError(err)

	response := dtos.Response{
		Code:     statusCode,
		Data:     data,
		Error:    standardErr,
		Message:  message,
		Metadata: metadata,
	}

	c.JSON(statusCode, response)

	return response
}

func ResponseError(c *gin.Context, err error) dtos.Response {
	var data any

	message := "failed"

	return response(c, e.GetStatusCode(err), data, nil, err, message)
}

func ResponseSuccess(c *gin.Context, data any, metadata any) dtos.Response {
	var statusCode int = http.StatusOK
	var err error = nil

	message := "success"

	return response(c, statusCode, data, metadata, err, message)
}
