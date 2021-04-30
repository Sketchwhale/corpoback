package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
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

type CompanyWithOwnership struct {
	C Company
	Owners []Ownership
	Subsidiaries []Ownership
}

type Ownership struct {
	Owner Company
	Owned Company
}

type TotalOverview struct {
	AllCompanies []Company
	AllOwnerships []Ownership
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
	overview := TotalOverview{AllCompanies: Companies, AllOwnerships: Ownerships}
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
				if company == ownership.Owner {
					subs = append(subs, ownership)
					//json.NewEncoder(w).Encode(ownership)
				}
				if company == ownership.Owned {
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
	fmt.Println(reqBody)
	var company Company
	json.Unmarshal(reqBody, &company)
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
	Owner string
	Owned string
}

func addBeneficialOwner (w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var newOwnershipJson NewOwnershipJson
	json.Unmarshal(reqBody, &newOwnershipJson)
	json.NewEncoder(w).Encode(newOwnershipJson)

	var masterCompany Company
	var slaveCompany Company

	foundOwner := false
	foundOwned := false
	for _, company := range Companies {
		if company.Id == newOwnershipJson.Owner {
			masterCompany = company
			foundOwner = true
		} else if company.Id == newOwnershipJson.Owned {
			slaveCompany = company
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


	var newOwnership = Ownership{Owner: masterCompany, Owned: slaveCompany}
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
	log.Fatal(http.ListenAndServe(":10000", router))
}

func main () {
	Companies = []Company {
		Company{Id: "SonyEurope", Name: "Sony", Address:	"KonigsStrasse", City: "Berlin"},
		Company{Id: "SonyJapan", Name: "Sony", Address:	"Shinjuku 1-2-4", City: "Tokyo"},
		Company{Id: "Toshiba", Name: "Toshiba", Address:	"Ginza 2-5-6", City: "Tokyo"},
	}
	Ownerships =[]Ownership {
		Ownership{Owner: Companies[1], Owned: Companies[0]},
		Ownership{Owner: Companies[1], Owned: Companies[2]},
	}
	handleRequests()
}
