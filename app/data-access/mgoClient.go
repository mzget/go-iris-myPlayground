package database

import (
	"gowork/app/utils"

	"github.com/globalsign/mgo"

	"log"
)

// MgoConnect use for init mongodb connection.
func MgoConnect(configuration utils.Configuration) *mgo.Session {
	// Database connection
	var connection = ""
	if configuration.Env == "Staging" {
		connection = configuration.DbDev
	} else {
		connection = configuration.DbPD
	}
	session, err := mgo.Dial(connection)
	if nil != err {
		log.Panic(err.Error())
	}
	session.SetMode(mgo.Monotonic, true)

	return session
}
