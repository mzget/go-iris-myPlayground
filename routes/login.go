package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"log"
	"time"
)

// MySigningKey key for jwt secret.
const MySigningKey string = "MySecret1234"

// MyCustomClaims for claim jwt payload.
type MyCustomClaims struct {
	username string
	password string
	_id      string
	jwt.StandardClaims
}

// Login with username, password ...
func Login(ctx iris.Context) {
	username, password := ctx.PostValue("username"), ctx.PostValue("password")

	// Create the Claims
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	claims := MyCustomClaims{
		username,
		password,
		"2",
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
	ctx.JSON(iris.Map{"data": ss})
}
