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
package main

import (
	"../database"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	districts map[string]citizen.Address
	genders   map[string]string = map[string]string{
		"1": "M",
		"2": "F",
	}
)

// loadDistrics reads the CSV file that contains the information
// related to registered voting districts from Costa Rica and
// store then into a memory map for fast consulting when
// reading the whole citizens database.
func loadDistrics() {
	districts = map[string]citizen.Address{}

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

		address := &citizen.Address{province, city, district}
		districts[id] = *address
	}
}

// loadPeople reads the CSV file containing the complete
// database of Costa Rican registered voters into a
// database aimed to interact with a web API server.
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

		person := &citizen.Citizen{id, name, first, last, gender, address}
		err = person.Create()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}

// getRecord transforms the given "record" which its just a text
// and returns the equivalent of it but capitalized and without
// trailing spaces
func getRecord(record string) string {
	return strings.TrimSpace(strings.Title(strings.ToLower(record)))
}

func main() {
	err := citizen.CreateIndex()
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	loadDistrics()
	loadPeople()
}
