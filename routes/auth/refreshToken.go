package auth

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"

	"gowork/app/utils"

	"log"
	"time"
)

// RefreshToken with old token and then return new token to client.
func RefreshToken(ctx iris.Context) {
	user := ctx.Values().Get("jwt").(*jwt.Token)

	my := utils.MyCustomClaims{}
	bytes, _ := json.Marshal(user.Claims)
	if err := json.Unmarshal(bytes, &my); err != nil {
		log.Print(err.Error())

		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"message": err.Error()})
		return
	}

	// Create the Claims
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	claims := utils.MyCustomClaims{
		my.Username,
		my.Password,
		my.ID,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "nattapon.r@live.com",
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString([]byte(utils.MySigningKey))

	ctx.JSON(iris.Map{"data": ss})
}
