package auth

import (
	"github.com/kataras/iris"
	// "log"

	"gowork/app/security"
	"gowork/app/utils"
)

// VerifyToken use to serve verify middleware.
func VerifyToken(ctx iris.Context, jwtHandler *middleware.Middleware) error {
	token, tokenError := jwtHandler.ValidationToken(ctx)
	if tokenError != nil {
		utils.ResponseFailure(ctx, iris.StatusNonAuthoritativeInfo, "", tokenError.Error())
		return tokenError
	}
	parseError := jwtHandler.ParseToken(ctx, token, &utils.MyCustomClaims{})
	if parseError != nil {
		utils.ResponseFailure(ctx, iris.StatusNonAuthoritativeInfo, "", parseError.Error())
		return parseError
	}

	return nil
}
