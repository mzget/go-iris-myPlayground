package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"

	"github.com/globalsign/mgo/bson"

	"gowork/app/data-access"
	"gowork/app/utils"

	"log"
	"time"
)

// MySigningKey key for jwt secret.
const MySigningKey string = "MySecret1234"

// MyCustomClaims for claim jwt payload.
type MyCustomClaims struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	ID       string `json:"id,omitempty"`
	jwt.StandardClaims
}

// Login with username, password ...
func Login(ctx iris.Context) {
	email, password := ctx.PostValue("email"), ctx.PostValue("password")
	user := User{
		Email:    email,
		Password: password,
	}
	err := user.Validate()
	if err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, nil, err)
		return
	}

	c := ctx.Values().Get("config")
	config, _ := c.(utils.Configuration)

	session := database.GetMgoSession()
	coll := session.DB(config.DbName).C(config.UserCollection)
	if err := coll.Find(bson.M{"email": email, "password": password}).One(&user); err != nil {
		utils.ResponseFailure(ctx, iris.StatusBadRequest, "NoContent", err)
		return
	}

	// Create the Claims
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	claims := MyCustomClaims{
		user.Email,
		user.Password,
		user.ID.String(),
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "nattapon.r@live.com",
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(MySigningKey))
	if err != nil {
		log.Panic(err)
	}

	utils.ResponseSuccess(ctx, ss)
}
