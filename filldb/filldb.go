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
