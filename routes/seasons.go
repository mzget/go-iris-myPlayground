package routes

import (
	"github.com/kataras/iris"

	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"gowork/app/data-access"
	"log"
)

type Season struct {
	no        int
	name      string
	programId string
}

// Query all Seasons.
func Seasons(ctx iris.Context) {
	client := database.Client

	coll := client.Database(database.CARTOON).Collection(database.SEASONS)
	cursor, err := coll.Find(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var results []string
	for cursor.Next(context.Background()) {
		elem := bson.NewDocument()
		if err := cursor.Decode(elem); err != nil {
			log.Fatal("Err", err)
		}

		/*
			bytes, err := elem.MarshalBSON()
			if err != nil {
				log.Fatal(err)
			}

			doc, err := bson.ToExtJSON(false, bytes)
			if err != nil {
				log.Fatal(err)
			}
			println(doc)
		*/

		text := elem.ToExtJSON(false)
		// do something with elem....
		results = append(results, text)
	}

	ctx.JSON(results)
}
