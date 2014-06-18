package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"labix.org/v2/mgo"
	"os"
	"strings"
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
	districts  map[string]Address
	genders    map[string]string = map[string]string{
		"1": "M",
		"2": "F",
	}
)

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		dialInfo := mgo.DialInfo{}
		dialInfo.Addrs = []string{os.Getenv("NT_DB_ADDRESS")}
		dialInfo.Username = os.Getenv("NT_DB_USER")
		dialInfo.Password = os.Getenv("NT_DB_PASS")
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

	c := session.DB(dbname).C("persons")
	err = c.Insert(&person)
	return
}

func createIndex() (err error) {
	// Create MongoDB session
	session := getSession()
	defer session.Close()

	c := session.DB(dbname).C("persons")

	// Index
	index := mgo.Index{
		Key:        []string{"identification"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = c.EnsureIndex(index)
	return
}

func loadDistrics() {
	districts = map[string]Address{}

	file, err := os.Open("Distelec.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		id := record[0]
		province := getRecord(record[1])
		city := getRecord(record[2])
		district := getRecord(record[3])

		address := &Address{province, city, district}
		districts[id] = *address
	}
}

func loadPeople() {
	file, err := os.Open("PADRON_COMPLETO.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		id := record[0]
		districtId := record[1]
		gender := genders[record[2]]
		name := getRecord(record[5])
		first := getRecord(record[6])
		last := getRecord(record[7])
		address := districts[districtId]

		person := &Person{id, name, first, last, gender, address}
		err = person.Create()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}

func getRecord(record string) string {
	return strings.TrimSpace(strings.Title(strings.ToLower(record)))
}

func main() {
	err := createIndex()
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	loadDistrics()
	loadPeople()
}
