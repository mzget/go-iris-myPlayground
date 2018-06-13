package user

import (
	// "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"

	"github.com/globalsign/mgo/bson"

	"gowork/app/data-access"
	"gowork/app/utils"

	"gowork/models"

	// "log"
	// "time"
	"encoding/json"
)

// GetUser will return user data.
func GetUser(ctx iris.Context) {
	token := ctx.Values().Get("user")

	myClaim := utils.MyCustomClaims{}
	bytes, _ := json.Marshal(token)
	if err := json.Unmarshal(bytes, &myClaim); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, err.Error())
		return
	}

	config := utils.ConfigParser(ctx)

	user := models.User{}

	session := database.GetMgoSession()
	coll := session.DB(config.DbName).C(config.UserCollection)
	if err := coll.Find(bson.M{"_id": bson.ObjectIdHex(myClaim.ID)}).Select(bson.M{"password": 0}).One(&user); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, "NoContent", err)
		return
	}

	utils.ResponseSuccess(ctx, user)
}

// PostUser use to update user data.
func PostUser(ctx iris.Context) {
	myClaim, err := utils.TokenParser(ctx)
	if err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, err.Error())
	}
	config := utils.ConfigParser(ctx)

	user := models.User{}
	if err := ctx.ReadJSON(&user); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, err.Error())
		return
	}

	m := bson.M{}
	if user.Firstname != "" {
		m["firstname"] = user.Firstname
	}
	if user.Lastname != "" {
		m["lastname"] = user.Lastname
	}
	if user.Gender != "" {
		m["gender"] = user.Gender
	}

	session := database.GetMgoSession()
	coll := session.DB(config.DbName).C(config.UserCollection)

	colQuerier := bson.M{"_id": bson.ObjectIdHex(myClaim.ID)}
	change := bson.M{"$set": m, "$currentDate": bson.M{"lastModified": true}}
	if err = coll.Update(colQuerier, change); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, err.Error())
		return
	}

	if err := coll.Find(bson.M{"_id": bson.ObjectIdHex(myClaim.ID)}).Select(bson.M{"password": 0}).One(&user); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, user)
}
