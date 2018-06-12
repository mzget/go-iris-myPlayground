package user

import (
	// "encoding/json"
	"github.com/kataras/iris"
	// "log"
	// "time"
	"gowork/app/data-access"
	"gowork/app/utils"
)

// User model.
type User struct {
	Name     string
	Password string
	Email    string
}

// Register user.
func Register(ctx iris.Context) {
	email, password := ctx.PostValue("email"), ctx.PostValue("password")
	c := ctx.Values().Get("config")
	config, _ := c.(utils.Configuration)

	var user = User{
		Email:    email,
		Password: password,
	}

	var session = database.GetMgoSession()
	coll := session.DB(config.DbName).C(config.UserCollection)
	if err := coll.Insert(user); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"message": err})

		return
	}

	ctx.JSON(iris.Map{"data": map[string]bool{"success": true}})
}
