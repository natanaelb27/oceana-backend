package utils

import (
	"net/http"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, statusCode int, data any, metadata any, err any, message string) dtos.Response {
	response := dtos.Response{
		Code:     statusCode,
		Data:     data,
		Error:    err,
		Message:  message,
		Metadata: metadata,
	}

	c.JSON(statusCode, response)

	return response
}

func ResponseError(c *gin.Context, err error) dtos.Response {
	var data any
	var message string
	var newErr string
	if err != nil {
		message = TranslateError(err)
		if message == "-" {
			message = err.Error()
			newErr = http.StatusText(e.GetStatusCode(err))
		} else {
			newErr = http.StatusText(http.StatusBadRequest)
		}
	}

	return Response(c, e.GetStatusCode(err), data, nil, newErr, message)
}

func ResponseSuccess(c *gin.Context, data any, metadata any) dtos.Response {
	var statusCode int = http.StatusOK

	message := "success"

	return Response(c, statusCode, data, metadata, nil, message)
}
