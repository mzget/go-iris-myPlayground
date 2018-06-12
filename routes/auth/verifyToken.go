package auth

import (
	"github.com/kataras/iris"
	// "log"

	"gowork/app/security"
	"gowork/app/utils"
)

// VerifyToken use to serve verify middleware.
func VerifyToken(ctx iris.Context, jwtHandler *jwtmiddleware.Middleware) {
	token, tokenError := jwtHandler.ValidationToken(ctx)
	if tokenError != nil {
		utils.ResponseFailure(ctx, iris.StatusNonAuthoritativeInfo, "", tokenError)
		return
	}
	parseError := jwtHandler.ParseToken(ctx, token, &MyCustomClaims{})
	if parseError != nil {
		utils.ResponseFailure(ctx, iris.StatusNonAuthoritativeInfo, "", parseError)
		return
	}

	ctx.Next()
}
