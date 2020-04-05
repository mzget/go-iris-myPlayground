package utils

import (
	"github.com/kataras/iris"
)

// ResponseSuccess when everything ok.
func ResponseSuccess(ctx iris.Context, data interface{}) {
	ctx.JSON(iris.Map{"data": data})
}

// ResponseFailure when everything so sad.
func ResponseFailure(ctx iris.Context, statusCode int, message interface{}, err interface{}) {
	ctx.StatusCode(statusCode)
	ctx.JSON(iris.Map{"message": message, "error": err})
}
