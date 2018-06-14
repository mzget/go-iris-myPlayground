// Package tollboothic v2(latest) provides rate-limiting logic to iris request handlers.
package middleware

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/errors"
	"github.com/didip/tollbooth/limiter"
	"github.com/kataras/iris/context"
)

// LimitHandler is a middleware that performs
// rate-limiting given a "limiter" configuration.
//
// Read more at: https://github.com/didip/tollbooth
// And https://github.com/didip/tollbooth_iris
func LimitHandler(ctx context.Context, l *limiter.Limiter) *errors.HTTPError {
	httpError := tollbooth.LimitByRequest(l, ctx.ResponseWriter(), ctx.Request())
	return httpError
}
