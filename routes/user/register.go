package user

import (
	// "encoding/json"
	"github.com/kataras/iris"

	"github.com/globalsign/mgo/bson"
	// "log"
	// "time"
	"gowork/app/data-access"
	"gowork/app/utils"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// User model.
type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string
	Password string
	Email    string
}

func (a User) Validate() error {
	return validation.ValidateStruct(&a,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Email, validation.Required, is.Email),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Password, validation.Required, validation.Length(8, 32)),
	)
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
	err := user.Validate()
	if err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, "", err)
		return
	}

	var session = database.GetMgoSession()
	coll := session.DB(config.DbName).C(config.UserCollection)
	if err := coll.Insert(user); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, "", err)
		return
	}

	utils.ResponseSuccess(ctx, map[string]bool{"success": true})
}
