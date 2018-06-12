package routes

import (
	"github.com/kataras/iris"
	"log"

	"gowork/app/security"
)

// VerifyToken use to serve verify middleware.
func VerifyToken(ctx iris.Context, jwtHandler *jwtmiddleware.Middleware) {
	token, tokenError := jwtHandler.ValidationToken(ctx)
	if tokenError != nil {
		log.Print(tokenError.Error())
		ctx.JSON(iris.Map{"message": tokenError.Error()})
		return
	}
	parseError := jwtHandler.ParseToken(ctx, token, &MyCustomClaims{})
	if parseError != nil {
		ctx.JSON(iris.Map{"message": parseError.Error()})
		return
	}

	ctx.Next()
}
