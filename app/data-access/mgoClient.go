package database

import (
	"gowork/app/utils"

	// "fmt"

	"github.com/globalsign/mgo"
)

var mgoSession *mgo.Session

// MgoConnect use for init mongodb connection.
func MgoConnect(configuration utils.Configuration) *mgo.Session {
	// fmt.Println(configuration)

	// Database connection
	var connection = ""
	if configuration.Env == "Staging" {
		connection = configuration.DbDev
	} else {
		connection = configuration.DbPD
	}
	session, err := mgo.Dial(connection)
	if nil != err {
		panic(err.Error())
	}
	session.SetMode(mgo.Monotonic, true)

	mgoSession = session

	return session
}

// GetMgoSession use for get session object.
func GetMgoSession() *mgo.Session {
	return mgoSession
}
