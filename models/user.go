package models

import (
	"github.com/globalsign/mgo/bson"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"time"
)

// User model.
type User struct {
	ID           bson.ObjectId `json:"id" bson:"_id,omitempty" `
	Name         string        `json:"name"`
	Firstname    string        `json:"firstname"`
	Lastname     string        `json:"lastname"`
	Gender       string        `json:"gender"`
	Birthdate    string        `json:"birthdate"`
	Password     string        `json:"password"`
	Email        string        `json:"email"`
	CreateAt     time.Time     `json:"createAt" bson:"createAt,omitempty"`
	LastModified time.Time     `json:"lastModified" bson:"lastModified,omitempty"`
}

func (a User) Validate() error {
	return validation.ValidateStruct(&a,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Email, validation.Required, is.Email),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Password, validation.Required, validation.Length(8, 32)),
	)
}
