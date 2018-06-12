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
			return []byte(routes.MySigningKey), nil
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
		routes.VerifyToken(ctx, jwtHandler)
	})
	apiRoutes.Get("/", func(ctx iris.Context) {
		log.Print(ctx.GetHeader(utils.ApiVersion))
		// ctx.HTML("<h1>Welcome</h1>")
		user := ctx.Values().Get("jwt").(*jwt.Token)
		log.Println(user.Claims)
		ctx.JSON(user)
	})
	apiRoutes.Get("/refreshToken", routes.RefreshToken)

	var authRoutes = app.Party("/auth")
	authRoutes.Post("/login", routes.Login)
	authRoutes.Post("/register", user.Register)

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

	// Method:   GET
	// Resource: http://localhost:8080/hello
	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})

	// registers a custom handler for 404 not found http (error) status code,
	// fires when route not found or manually by ctx.StatusCode(iris.StatusNotFound).
	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)
	// var sb strings.Builder
	// sb.WriteString(":")
	// sb.WriteString(configuration.Port)
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
