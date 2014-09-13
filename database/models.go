// The MIT License (MIT)
//
// Copyright (c) 2014 Greivin López
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package citizen

import (
	"github.com/greivinlopez/skue"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

// Citizen represents a Costa Rican citizen.
type Citizen struct {
	Identification string
	FirstName      string
	LastName       string
	SurName        string
	Gender         string
	Address        Address
}

// Address will contain the registered voting district of the citizen.
type Address struct {
	Province string
	City     string
	District string
}

var (
	mgoSession *mgo.Session
	dbname     string = "people"
	collection string = "citizens"
)

// ----------------------------------------------------------------------------
// 			MONGODB
// ----------------------------------------------------------------------------
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

// CreateIndex creates an index associated to the id
// of the citizen on the MongoDB collection in order
// to make faster queries when searching for citizen
// information.
func CreateIndex() (err error) {
	// Create MongoDB session
	session := getSession()
	defer session.Close()
	// Get the collection
	c := session.DB(dbname).C(collection)
	// Create the index
	index := mgo.Index{
		Key:      []string{"identification"},
		Unique:   true,
		DropDups: true,
	}
	err = c.EnsureIndex(index)
	return
}

// ----------------------------------------------------------------------------

// New creates a new empty Citizen with the provided identification
// All the other fields will be empty at first.
func New(id string) *Citizen {
	return &Citizen{
		Identification: id,
		FirstName:      "",
		LastName:       "",
		SurName:        "",
		Gender:         "",
		Address:        Address{},
	}
}

// ----------------------------------------------------------------------------
// 			skue.DatabasePersistor implementation
// ----------------------------------------------------------------------------

func (citizen *Citizen) Read(cache skue.MemoryCacher) (err error) {
	// Create MongoDB session
	session := getSession()
	defer session.Close()

	c := session.DB(dbname).C(collection)
	query := bson.M{"identification": citizen.Identification}
	err = c.Find(query).One(&citizen)
	return err
}

func (citizen *Citizen) Create() (err error) {
	// Create MongoDB session
	session := getSession()
	defer session.Close()

	c := session.DB(dbname).C(collection)
	err = c.Insert(&citizen)
	return
}

func (citizen *Citizen) Update(cache skue.MemoryCacher) (err error) {
	return nil
}

func (citizen *Citizen) Delete(cache skue.MemoryCacher) (err error) {
	return nil
}

// ----------------------------------------------------------------------------
