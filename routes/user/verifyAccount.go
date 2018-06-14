package user

import (
	"time"
	// "encoding/json"
	"github.com/kataras/iris"

	"github.com/globalsign/mgo/bson"
	// "time"
	"fmt"
	"gowork/app/controller"
	"gowork/app/data-access"
	"gowork/app/models"
	"gowork/app/utils"
)

// VerifyAccount use for email activation process.
func VerifyAccount(ctx iris.Context) {
	secret := ctx.PostValue("secret")
	config := utils.ConfigParser(ctx)

	if secret == "" {
		fmt.Println("Missing secret data")
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, "Missing secret data")
		return
	}

	email := controller.Decrypt(config.GeneratedLinkKey, secret)

	var session = database.GetMgoSession()
	coll := session.DB(config.DbName).C(config.UserCollection)
	filter := bson.M{"email": email}
	change := bson.M{"$set": bson.M{"verified": true, "verifiedAt": time.Now()},
		"$currentDate": bson.M{"lastModified": true}}
	if err := coll.Update(filter, change); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, err.Error())
		return
	}

	var user = models.User{}
	if err := coll.Find(filter).One(&user); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, user)
}

// ResendActivationEmail use for send email to user again and again.
func ResendActivationEmail(ctx iris.Context) {
	email := ctx.PostValue("email")
	if email == "" {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, "Missing email field")
		return
	}

	user := models.User{}

	config := utils.ConfigParser(ctx)
	session := database.GetMgoSession()
	coll := session.DB(config.DbName).C(config.UserCollection)
	filter := bson.M{"email": email}
	if err := coll.Find(filter).One(&user); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, err.Error())
		return
	}

	// encrypt value to base64
	cryptoText := controller.Encrypt(config.GeneratedLinkKey, user.Email)

	utils.ResponseSuccess(ctx, iris.Map{
		"success": true,
		"message": "Verification email will send to you as " + user.Email,
		"secret":  cryptoText,
	})
}
