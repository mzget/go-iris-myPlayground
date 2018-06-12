package user

import (
	"encoding/json"
	"github.com/kataras/iris"
	"log"
	"time"
)

// User model.
type User struct {
	Name     string
	Username string
	Password string
}

// Register user.
func Register(ctx iris.Context) {
	var result = User{}

	// Database name and collection name
	// car-db is database name car is collation name
	c := session.DB("car-db").C("car")
	c.Insert(&Car{"Audi", "Luxury car"})

	ctx.JSON(result)
}
