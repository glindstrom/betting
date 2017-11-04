package db

import (
	"gopkg.in/mgo.v2"
	"os"
)

var (
	// Session stores mongo session
	Session *mgo.Session
)

const (
	// MongoDBUrl is the default mongodb url that will be used to connect to the
	// database.
	MongoDBUrl = "mongodb://localhost"
	DataBase   = "betting_db"
)

// Connect connects to mongodb
func Connect() {
	uri := os.Getenv("MONGODB_URL")

	if len(uri) == 0 {
		uri = MongoDBUrl
	}

	Session = getSession()
}

func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial(MongoDBUrl)

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}
