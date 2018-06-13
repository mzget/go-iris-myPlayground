package utils

import (
	"github.com/kataras/iris"

	"encoding/json"
)

// TokenParser use for parse token object to Claim interface.
func TokenParser(ctx iris.Context) (MyCustomClaims, error) {
	token := ctx.Values().Get("user")

	myClaim := MyCustomClaims{}
	bytes, _ := json.Marshal(token)
	if err := json.Unmarshal(bytes, &myClaim); err != nil {
		return myClaim, err
	}

	return myClaim, nil
}
