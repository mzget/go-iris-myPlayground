package user

import (
	// "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"

	"github.com/globalsign/mgo/bson"

	"gowork/app/data-access"
	"gowork/app/utils"

	"gowork/routes/auth"

	"gowork/models"

	// "log"
	// "time"
	"encoding/json"
)

// GetUser will return user data.
func GetUser(ctx iris.Context) {
	token := ctx.Values().Get("user")

	myClaim := auth.MyCustomClaims{}
	bytes, _ := json.Marshal(token)
	if err := json.Unmarshal(bytes, &myClaim); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, err.Error())
		return
	}

	c := ctx.Values().Get("config")
	config, _ := c.(utils.Configuration)

	user := models.User{}

	session := database.GetMgoSession()
	coll := session.DB(config.DbName).C(config.UserCollection)
	if err := coll.Find(bson.M{"_id": bson.ObjectIdHex(myClaim.ID)}).Select(bson.M{"password": 0}).One(&user); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, "NoContent", err)
		return
	}

	utils.ResponseSuccess(ctx, user)
}
