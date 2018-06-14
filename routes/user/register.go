package user

import (
	"fmt"
	"time"
	// "encoding/json"
	"github.com/kataras/iris"

	"github.com/globalsign/mgo/bson"
	// "log"
	// "time"
	"gowork/app/controller"
	"gowork/app/data-access"
	"gowork/app/models"
	"gowork/app/utils"
)

// Register user.
func Register(ctx iris.Context) {
	email, password := ctx.PostValue("email"), ctx.PostValue("password")
	c := ctx.Values().Get("config")
	config, _ := c.(utils.Configuration)

	var user = models.User{
		Email:    email,
		Password: password,
		CreateAt: time.Now(),
	}
	err := user.Validate()
	if err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, "", err)
		return
	}

	var session = database.GetMgoSession()
	coll := session.DB(config.DbName).C(config.UserCollection)
	// Find email already register first.
	num, notFound := coll.Find(bson.M{"email": user.Email}).Count()

	if notFound != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, notFound)
		return
	}

	if num > 0 {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, "Email already used.", nil)
		return
	}

	if err := coll.Insert(user); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, "", err)
		return
	}

	utils.ResponseSuccess(ctx, iris.Map{
		"success": true,
		"message": "Verification email will send to you as " + user.Email})

	autoSendEmail(ctx, user.Email)
}

func autoSendEmail(ctx iris.Context, email string) {
	config := utils.ConfigParser(ctx)

	key := []byte(config.GeneratedLinkKey)
	// encrypt value to base64
	cryptoText := controller.Encrypt(key, email)

	fmt.Println(cryptoText)
}
