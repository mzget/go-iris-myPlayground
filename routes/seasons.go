package seasons

import (
	"github.com/kataras/iris"

	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"gowork/app"
	"log"
)

// Query all Seasons.
func Seasons(ctx iris.Context) {
	client := database.Client

	coll := client.Database(database.CARTOON).Collection(database.SEASONS)
	cursor, err := coll.Find(context.Background(), bson.NewDocument())
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

		// do something with elem....
		results = append(results, elem.ToExtJSON(true))
	}

	ctx.JSON(results)
}
