package main

import (
	"../database"
	"github.com/greivinlopez/skue"
	"gopkg.in/martini.v1"
	"net/http"
	"os"
)

var apiKey string

// ----------------------------------------------------------------------------
// HANDLERS

func getPersonHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
	id := params["id"]
	person, err := database.ReadPerson(id)
	if err != nil {
		skue.ServiceResponse(w, http.StatusNotFound, "Citizen not found")
	} else {
		skue.ToJson(w, http.StatusOK, person)
	}
}

func main() {
	apiKey = os.Getenv("CZ_API_KEY")

	m := martini.Classic()

	// Validate an API key: Authorization
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-API-KEY") != apiKey {
			skue.ServiceResponse(res, http.StatusUnauthorized, "You are not authorized to access this resource.")
		}
	})

	// Citizens API
	m.Get("/citizens/:id", getPersonHandler)
	m.Any("/citizens/:id", skue.NotAllowed)

	// Running on an unassigned port by IANA: http://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers
	http.ListenAndServe(":3020", m)
}
