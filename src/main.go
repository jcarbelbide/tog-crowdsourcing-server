package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

type WorldInformation struct {
	WorldNumber int    `json:"world_number"`
	Hits        int    `json:"hits"`
	StreamOrder string `json:"stream_order"`
}

var (
	db                *sql.DB
	lastResetTimeUnix int64
)

const (
	minWorldValue         = 300
	maxWorldValue         = 1000 // allowing for Jagex to add plenty more servers.
	pollJS5ServerInterval = 5000
)

func main() {
	// TODO: Change prints to logging

	// Init Logger
	logFile := initLogging()
	defer logFile.Close()

	// Start Polling JS5 ServerInfo
	go pollJS5Server()

	// Database
	initDatabase()

	// Init Router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	r.HandleFunc("/worldinformation", getWorldInformation).Methods("GET")
	r.HandleFunc("/worldinformation", postWorldInformation).Methods("POST")

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

	remoteIPAddress := getRemoteIPAddressFromRequest(r)
	var newWorldInformation WorldInformation
	json.NewDecoder(r.Body).Decode(&newWorldInformation)

	dataIsValid := verifyDataIsValid(newWorldInformation)
	if !dataIsValid {
		return
	} // return early if data is garbage

	ipAddressIsValid, err := verifyIPAddressIsValid(newWorldInformation, remoteIPAddress, db)
	if err != nil {
		err = createAndLogCustomError(err, "Error when trying to verify if IP address was valid")
	}

	if ipAddressIsValid {
		addIPAndWorldToDB(newWorldInformation, remoteIPAddress, db)
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
// -------------------------- Poll JS5 ServerInfo --------------------------- //
// -------------------------------------------------------------------------- //
func pollJS5Server() {
	var liveLastResetTimeUnix int64 = 0
	initLastResetTimeUnix()

	for { // Poll forever in goroutine
		getJS5ServerInfo(&liveLastResetTimeUnix)
		time.Sleep(pollJS5ServerInterval * time.Millisecond)

		if liveLastResetTimeUnix != lastResetTimeUnix {
			clearHitsOnServerReset(db)
			lastResetTimeUnix = liveLastResetTimeUnix
		}
	}
}

func initLastResetTimeUnix() {
	getJS5ServerInfo(&lastResetTimeUnix)
}
