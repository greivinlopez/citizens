package database

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

type Person struct {
	Identification string
	FirstName      string
	LastName       string
	SurName        string
	Gender         string
	Address        Address
}

type Address struct {
	Province string
	City     string
	District string
}

var (
	mgoSession *mgo.Session
	dbname     string = "people"
	collection string = "persons"
)

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		dialInfo := mgo.DialInfo{}
		dialInfo.Addrs = []string{os.Getenv("CZ_DB_ADDRESS")}
		dialInfo.Username = os.Getenv("CZ_DB_USER")
		dialInfo.Password = os.Getenv("CZ_DB_PASS")
		dialInfo.Database = dbname
		mgoSession, err = mgo.DialWithInfo(&dialInfo)
		if err != nil {
			panic(err) // no, not really
		}
	}
	return mgoSession.Clone()
}

func (person *Person) Create() (err error) {
	// Create MongoDB session
	session := getSession()
	defer session.Close()

	c := session.DB(dbname).C(collection)
	err = c.Insert(&person)
	return
}

func ReadPerson(id string) (person Person, err error) {
	// Create MongoDB session
	session := getSession()
	defer session.Close()

	c := session.DB(dbname).C(collection)
	query := bson.M{"identification": id}
	person = Person{}
	err = c.Find(query).One(&person)
	return person, err
}

func CreateIndex() (err error) {
	// Create MongoDB session
	session := getSession()
	defer session.Close()

	c := session.DB(dbname).C(collection)

	// Index
	index := mgo.Index{
		Key:      []string{"identification"},
		Unique:   true,
		DropDups: true,
	}

	err = c.EnsureIndex(index)
	return
}
