package main

import (
	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
	"github.com/dgrijalva/jwt-go"

	"github.com/kataras/iris"
	"github.com/kataras/iris/core/host"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"

	"gowork/app/data-access"
	"gowork/app/security"
	"gowork/app/utils"

	"gowork/routes"
	"gowork/routes/auth"
	"gowork/routes/user"

	"log"
	"os"
	"path"
	"path/filepath"
	// "strings"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// first parameter is the request path
	// second is the system directory
	app.StaticWeb("/static", "./assets")

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic(err)
	}
	configuration := utils.GetConfig(path.Join(dir, "conf.json"))
	if configuration.Env == "Staging" {
		yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
			On:       true,
			DocTitle: "Iris",
			DocPath:  "./assets/apidoc.html",
			BaseUrls: map[string]string{"Production": "", "Staging": ""},
		})
		app.Use(irisyaag.New()) // <- IMPORTANT, register the middleware.
	}

	// Database connection
	session := database.MgoConnect(configuration)
	defer session.Close()

	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(utils.MySigningKey), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
		Expiration:    true,
	})

	app.Use(func(ctx iris.Context) {
		ctx.Values().Set("config", configuration)
		ctx.Next()
	})
	var apiRoutes = app.Party("/api")
	apiRoutes.Use(func(ctx iris.Context) {
		err := auth.VerifyToken(ctx, jwtHandler)
		if err == nil {
			user := ctx.Values().Get("jwt").(*jwt.Token)
			ctx.Values().Set("user", user.Claims)
			ctx.Next()
		}
	})
	apiRoutes.Get("/refreshToken", auth.RefreshToken)
	apiRoutes.Get("/user", user.GetUser)
	apiRoutes.Post("/user", user.PostUser)

	var authRoutes = app.Party("/auth")
	authRoutes.Post("/login", user.Login)
	authRoutes.Post("/register", user.Register)
	authRoutes.Post("/verifyAccount", user.VerifyAccount)
	authRoutes.Post("/resendEmail", user.ResendActivationEmail)

	/* Official mongodb client.

		client := database.Connect()
	 	defer client.Disconnect(context.Background())

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
	*/
	app.Get("/seasons", routes.Seasons)

	// registers a custom handler for 404 not found http (error) status code,
	// fires when route not found or manually by ctx.StatusCode(iris.StatusNotFound).
	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)
	// var sb strings.Builder
	// sb.WriteString(":")
	// sb.WriteString(configuration.Port)
	app.Configure(iris.WithConfiguration(iris.Configuration{
		DisableAutoFireStatusCode: true,
	}))
	app.Run(iris.Addr(":"+configuration.Port, configureHost), iris.WithoutServerError(iris.ErrServerClosed))
}

func notFoundHandler(ctx iris.Context) {
	ctx.HTML("Custom route for 404 not found http code, here you can render a view, html, json <b>any valid response</b>.")
}

func configureHost(su *host.Supervisor) {
	// here we have full access to the host that will be created
	// inside the `app.Run` function or `NewHost`.
	//
	// we're registering a shutdown "event" callback here:
	su.RegisterOnShutdown(func() {
		println("server is closed")
	})
	// su.RegisterOnError
	// su.RegisterOnServe
}
