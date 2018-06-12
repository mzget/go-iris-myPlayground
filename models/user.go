package models

import (
	"github.com/globalsign/mgo/bson"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// User model.
type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string
	Password string
	Email    string
}

func (a User) Validate() error {
	return validation.ValidateStruct(&a,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Email, validation.Required, is.Email),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Password, validation.Required, validation.Length(8, 32)),
	)
}
