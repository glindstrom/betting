package db

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

var (
	// Session stores mongo session
	Session *mgo.Session
)

const (
	// MongoDBUrl is the default mongodb url that will be used to connect to the
	// database.
	MongoDBUrl      = "mongodb://localhost"
	DataBase        = "betting_db"
	GamesCollection = "games"
)

// database
var DB *mgo.Database

// collections
var Games *mgo.Collection

func init() {
	// get a mongo sessions
	//s, err := mgo.Dial("mongodb://bond:moneypenny007@localhost/bookstore")
	s, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
	}

	if err = s.Ping(); err != nil {
		panic(err)
	}

	DB = s.DB(DataBase)
	Games = DB.C(GamesCollection)

	index := mgo.Index{
		Key:        []string{"status"},
		Unique:     false,
		Background: true,
		Sparse:     true,
	}
	err = Games.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	index = mgo.Index{
		Key:        []string{"-dateTime"},
		Unique:     false,
		Background: true,
		Sparse:     true,
	}
	err = Games.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	index = mgo.Index{
		Key:        []string{"id538"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}
	err = Games.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	fmt.Println("You connected to your mongo database.")
}
