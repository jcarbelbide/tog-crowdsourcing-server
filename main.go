package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

type WorldInformation struct {
	WorldNumber int    `json:"world_number"`
	Hits        int    `json:"hits"`
	StreamOrder string `json:"stream_order"`
}

var db *sql.DB

const (
	MIN_WORLD_VALUE = 300
	MAX_WORLD_VALUE = 1000 // allowing for Jagex to add plenty more servers.
)

func main() {

	// Init Logger
	logFile := initLogging()
	defer logFile.Close()

	// Database
	db, _ = sql.Open("sqlite3", "./ToGWorldInformation.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS World_Information (wi_id INTEGER PRIMARY KEY AUTOINCREMENT, world_number INTEGER, hits INTEGER, stream_order TEXT)")
	statement.Exec()

	statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS IP_List (ip_id INTEGER PRIMARY KEY AUTOINCREMENT, world_number INTEGER, ip_address TEXT)")
	statement.Exec()

	// Init Router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	r.HandleFunc("/worldInformation", getWorldInformation).Methods("GET")
	r.HandleFunc("/worldInformation", postWorldInformation).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// -------------------------------------------------------------------------- //
// ---------------------------- Request Handlers ---------------------------- //
// -------------------------------------------------------------------------- //

func getWorldInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	worldInformationList := queryDBForAllWorldInformation(db)

	json.NewEncoder(w).Encode(worldInformationList)
	return
}

func postWorldInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newWorldInformation WorldInformation
	json.NewDecoder(r.Body).Decode(&newWorldInformation)

	dataIsValid := verifyDataIsValid(newWorldInformation)
	ipAddressIsValid, err := verifyIPAddressIsValid(newWorldInformation, r, db)
	if err != nil {
		log.Println(err)
	}

	if ipAddressIsValid {
		// TODO Add the IP + world to the database. Will just have 1 table for the IPs. Each row will be a world number + IP.

	} else {
		return
	}

	if dataIsValid {
		// First, check to see if the world + stream order combo is already in the db.
		entryExistsInDB, existingWorldInformation := queryDBForSpecificWorldInformation(newWorldInformation, db)

		// If it is, hits++ and update db
		if entryExistsInDB {
			incrementHitsOnExistingWorld(&existingWorldInformation, db)
			json.NewEncoder(w).Encode(existingWorldInformation)
		} else { // Else add it to the db
			newWorldInformation.Hits = 1
			addNewWorldInformation(newWorldInformation, db)
			json.NewEncoder(w).Encode(newWorldInformation)
		}

	}

}

// -------------------------------------------------------------------------- //
// ---------------------------- Helper Functions ---------------------------- //
// -------------------------------------------------------------------------- //

func createCustomError(err error, message string) error {
	return fmt.Errorf(message, err)
}

func initLogging() *os.File {
	file, err := os.OpenFile("./logs/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	return file
}
