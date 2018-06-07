package main

import (
	"github.com/kataras/iris"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"

	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
)

const CARTOON string = "test-webtoons"
const PROGRAMS string = "programs"

type Name struct {
	th string
	en string
}
type PROGRAM struct {
	name Name
}

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
	client, err := mongo.Connect(context.Background(), "mongodb://mzget:mzget1234@chitchats.ga:27017", nil)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer client.Disconnect(context.Background())

	// registers a custom handler for 404 not found http (error) status code,
	// fires when route not found or manually by ctx.StatusCode(iris.StatusNotFound).
	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)

	// Method:   GET
	// Resource: http://localhost:8080
	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.HTML("<h1>Welcome</h1>")
	})

	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		// Database name and collection name
		coll := client.Database(CARTOON).Collection(PROGRAMS)
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
