package user

import (
	"time"
	// "encoding/json"
	"github.com/kataras/iris"

	"fmt"
	"github.com/globalsign/mgo/bson"
	"log"
	// "time"
	"gowork/app/data-access"
	"gowork/app/utils"
	"gowork/models"

	"crypto/aes"
	"crypto/cipher"
	// "crypto/rand"
	"encoding/base64"
)

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		log.Panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprint(ciphertext)
}

// VerifyAccount use for user verification process.
func VerifyAccount(ctx iris.Context) {
	secret := ctx.PostValue("secret")
	config := utils.ConfigParser(ctx)

	key := []byte(config.GeneratedLinkKey)
	email := decrypt(key, secret)

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
