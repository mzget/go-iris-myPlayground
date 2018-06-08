package main

import (
	"github.com/kataras/iris"

	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/mongodb/mongo-go-driver/bson"
	"gowork/app/data-access"
	"gowork/app/security"
	"gowork/routes"
	"log"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// Database connection
	// session, err := mgo.Dial("mzget:mzget1234@chitchats.ga:27017")
	// if nil != err {
	// 	panic(err)
	// }
	// defer session.Close()
	// session.SetMode(mgo.Monotonic, true)

	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(routes.MySigningKey), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	apiRoutes := app.Party("/api")
	apiRoutes.Use(func(ctx iris.Context) {
		token, tokenError := jwtHandler.ValidationToken(ctx)
		if tokenError != nil {
			log.Print(tokenError.Error())
			ctx.JSON(iris.Map{"message": tokenError.Error()})
			return
		}
		parseError := jwtHandler.ParseToken(ctx, token, &routes.MyCustomClaims{})
		if parseError != nil {
			ctx.JSON(iris.Map{"message": parseError.Error()})
			return
		}

		ctx.Next()
	})

	client := database.Connect()
	defer client.Disconnect(context.Background())

	// registers a custom handler for 404 not found http (error) status code,
	// fires when route not found or manually by ctx.StatusCode(iris.StatusNotFound).
	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)

	// Method:   GET
	// Resource: http://localhost:8080
	apiRoutes.Get("/", func(ctx iris.Context) {
		// ctx.HTML("<h1>Welcome</h1>")

		user := ctx.Values().Get("jwt").(*jwt.Token)

		ctx.Writef("This is an authenticated request\n")
		ctx.Writef("Claim content:\n")
		ctx.JSON(user)
	})

	authRoutes := app.Party("/auth")
	authRoutes.Post("/login", routes.Login)

	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		// Database name and collection name
		coll := client.Database(database.CARTOON).Collection(database.PROGRAMS)
		cursor, err := coll.Find(context.Background(), bson.NewDocument())
		if err != nil {
			log.Fatal(err)
		}
		for cursor.Next(context.Background()) {
			elem := bson.NewDocument()
			if err := cursor.Decode(elem); err != nil {
				log.Fatal(err)
			}

			// do something with elem....
			log.Fatal("result", elem.ToExtJSON(true))
			ctx.WriteString(elem.ToExtJSON(true))
		}
	})
	app.Get("/seasons", routes.Seasons)

	// Method:   GET
	// Resource: http://localhost:8080/hello
	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})

	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func notFoundHandler(ctx iris.Context) {
	ctx.HTML("Custom route for 404 not found http code, here you can render a view, html, json <b>any valid response</b>.")
}
