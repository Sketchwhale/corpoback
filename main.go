package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"os"
)


type Company struct {
	Id string
	Name string
	Address string
	City string
	Country string
	Email string
	Phone string
}

func (c *Company) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Id *string
		Name *string
		Address *string
		City *string
		Country *string
	}{}
	complete := struct {
		Id string
		Name string
		Address string
		City string
		Country string
		Email string
		Phone string
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return err
	} else if (required.Id == nil || *required.Id == "") {
		err = fmt.Errorf ("Id cannot be empty or nil")
	} else if (required.Name == nil || *required.Name == "") {
		err = fmt.Errorf ("Name cannot be empty or nil")
	} else if (required.Address == nil || *required.Address == "") {
		err = fmt.Errorf ("Address cannot be empty or nil")
	} else if (required.City == nil || *required.City == "") {
		err = fmt.Errorf ("City cannot be empty or nil")
	} else if (required.Country == nil || *required.Country == "") {
		err = fmt.Errorf ("Country cannot be empty or nil")
	} else {
		err 		= json.Unmarshal(data, &complete)
		c.Id 		= complete.Id
		c.Name 		= complete.Name
		c.Address 	= complete.Address
		c.City 		= complete.City
		c.Country 	= complete.Country
		c.Email 	= complete.Email
		c.Phone 	= complete.Phone
	}
	return
}

type CompanyWithOwnership struct {
	C Company
	Owners []Ownership
	Subsidiaries []Ownership
}

type Ownership struct {
	OwnerId string
	OwnedId string
}

var Companies []Company
var Owners []Ownership
var Subsidiaries []Ownership
var Ownerships []Ownership

func returnAllCompanies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: returnAllCompanies")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	encodeAllCompanies(w)
}

func encodeAllCompanies (w http.ResponseWriter) {

	overview := struct {
		AllCompanies []Company
		AllOwnerships []Ownership
	}{
		AllCompanies: Companies,
		AllOwnerships: Ownerships,
	}
	json.NewEncoder(w).Encode(overview)

}

func getCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Println("Endpoint hit: getting company with ID: " + key)
	var c Company
	var subs []Ownership
	var owners []Ownership
	for _, company := range Companies {

		if company.Id == key {
			c = company

			for _, ownership := range Ownerships {
				if company.Id == ownership.OwnerId {
					subs = append(subs, ownership)
				}
				if company.Id == ownership.OwnedId {
					owners = append(owners, ownership)
				}
			}

			companyWithOwnership := CompanyWithOwnership {C: c, Owners: owners, Subsidiaries:  subs}
			json.NewEncoder(w).Encode(companyWithOwnership)

		}

	}
}

/*
func deleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, company := range Companies {
		if company.Id == id {
			Companies = append(Companies[:index], Companies[index+1:]...)
		}
	}
}
*/

func createCompany(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: creating company")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	reqBody, _ := ioutil.ReadAll(r.Body)
	//fmt.Println(reqBody)
	var company Company
	err := json.Unmarshal(reqBody, &company)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	Companies = append(Companies, company)
	json.NewEncoder(w).Encode(company)
}

func updateCompany(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept-Encoding, Accept, Content-Type")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" { return }

	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("Endpoint hit: Updating company with ID: " + id)

	reqBody, _ := ioutil.ReadAll(r.Body)
	var update Company
	json.Unmarshal(reqBody, &update)

	for index, company := range Companies {
		if company.Id == id {
			fmt.Println("found a match")
			Companies[index] = update
			json.NewEncoder(w).Encode(Companies[index])
		}
	}
}

type NewOwnershipJson struct {
	OwnerId string
	OwnedId string
}

func addBeneficialOwner (w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var newOwnership Ownership
	json.Unmarshal(reqBody, &newOwnership)
	json.NewEncoder(w).Encode(newOwnership)

	foundOwner := false
	foundOwned := false
	for _, company := range Companies {
		if company.Id == newOwnership.OwnerId {
			foundOwner = true
		} else if company.Id == newOwnership.OwnedId {
			foundOwned = true
		}
	}

	if foundOwner == false {
		json.NewEncoder(w).Encode("No owner found")
		return;
	}
	if foundOwned == false {
		json.NewEncoder(w).Encode("No owned found")
		return;
	}

	Ownerships = append(Ownerships, newOwnership)
	json.NewEncoder(w).Encode(newOwnership)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/all", returnAllCompanies)
	router.HandleFunc("/company", createCompany).Methods("POST")
	//router.HandleFunc("/company/{id}", deleteCompany).Methods("DELETE")
	router.HandleFunc("/company/{id}", updateCompany).Methods("PUT", "OPTIONS")
	router.HandleFunc("/company/{id}", getCompany)
	router.HandleFunc("/addOwnership", addBeneficialOwner).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func main () {
	Companies = []Company {
		Company{Id: "SonyEurope", Name: "Sony", Address:	"KonigsStrasse", City: "Berlin"},
		Company{Id: "SonyJapan", Name: "Sony", Address:	"Shinjuku 1-2-4", City: "Tokyo"},
		Company{Id: "Toshiba", Name: "Toshiba", Address:	"Ginza 2-5-6", City: "Tokyo"},
	}
	Ownerships =[]Ownership {
		Ownership{OwnerId: Companies[1].Id, OwnedId: Companies[0].Id},
		Ownership{OwnerId: Companies[1].Id, OwnedId: Companies[2].Id},
	}
	handleRequests()
}
