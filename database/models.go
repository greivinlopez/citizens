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
	"github.com/greivinlopez/skue/database"
	"gopkg.in/mgo.v2"
)

var (
	ServerAddress string // The address to reach the MongoDB server
	Username      string // The username to connect with the MongoDB server
	Password      string // The password of the MongoDB user
	Database      string // The name of the database to store the models
	mongo         *mongodb.MongoDBPersistor
	mgoSession    *mgo.Session
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

// Creates a MongoDB persistor to interact with the database
func CreateMongoPersistor() {
	mongo = mongodb.New(ServerAddress, Username, Password, Database)
}

// ----------------------------------------------------------------------------
// 			MONGODB
// ----------------------------------------------------------------------------
func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		dialInfo := mgo.DialInfo{}
		dialInfo.Addrs = []string{ServerAddress}
		dialInfo.Username = Username
		dialInfo.Password = Password
		dialInfo.Database = Database
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
	c := session.DB("people").C("citizens")
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
	err = mongo.Read(cache, &citizen, "citizens", "identification", citizen.Identification)
	return
}

func (citizen *Citizen) Create() (err error) {
	err = mongo.Create(&citizen, "citizens")
	return
}

func (citizen *Citizen) Update(cache skue.MemoryCacher) (err error) {
	return nil
}

func (citizen *Citizen) Delete(cache skue.MemoryCacher) (err error) {
	return nil
}

func (citizen *Citizen) List() (results interface{}, err error) {
	return nil, nil
}

// ----------------------------------------------------------------------------
