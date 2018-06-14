package user

import (
	"time"
	// "encoding/json"
	"github.com/kataras/iris"

	"github.com/globalsign/mgo/bson"
	// "time"
	"gowork/app/controller"
	"gowork/app/data-access"
	"gowork/app/models"
	"gowork/app/utils"
)

// VerifyAccount use for user verification process.
func VerifyAccount(ctx iris.Context) {
	secret := ctx.PostValue("secret")
	config := utils.ConfigParser(ctx)

	key := []byte(config.GeneratedLinkKey)
	email := controller.Decrypt(key, secret)

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
