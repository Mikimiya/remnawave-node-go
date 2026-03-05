package api

import (
	"github.com/gin-gonic/gin"

	apperrors "github.com/hteppl/remnawave-node-go/internal/errors"
)

func ErrorHandler(code string, c *gin.Context) {
	errDef, ok := apperrors.GetError(code)
	if !ok {
		errDef = apperrors.ERRORS[apperrors.CodeInternalServerError]
	}

	c.JSON(errDef.HTTPCode, NewErrorResponse(
		c.Request.URL.Path,
		errDef.Message,
		errDef.Code,
	))
}
