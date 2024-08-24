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
	cache             *Cache
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

	// Cache
	cache = NewCache(NewDBWrapper(db))

	// Init Router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	r.HandleFunc("/worldinfo", getWorldInformation).Methods("GET")
	r.HandleFunc("/worldinfo", postWorldInformation).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}

// -------------------------------------------------------------------------- //
// ---------------------------- Request Handlers ---------------------------- //
// -------------------------------------------------------------------------- //
func getWorldInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	worldInformationList := cache.Get()

	err := json.NewEncoder(w).Encode(worldInformationList)
	if err != nil {
		err = createAndLogCustomError(err, "Encode to world information list")
	}

	return
}

func postWorldInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	remoteIPAddress := getRemoteIPAddressFromRequest(r)
	var newWorldInformation WorldInformation
	err := json.NewDecoder(r.Body).Decode(&newWorldInformation)
	if err != nil {
		err = createAndLogCustomError(err, "Error when trying to decode body json in postWorldInformation")
		return
	}

	dataIsValid := verifyDataIsValid(newWorldInformation)
	if !dataIsValid {
		return
	} // return early if data is garbage

	ipAddressIsValid, err := verifyIPAddressIsValid(newWorldInformation, remoteIPAddress, db)
	if err != nil {
		err = createAndLogCustomError(err, "Error when trying to verify if IP address was valid")
		return
	}

	if !ipAddressIsValid {
		return
	}

	addIPAndWorldToDB(newWorldInformation, remoteIPAddress, db)

	// First, check to see if the world + stream order combo is already in the db.
	entryExistsInDB, existingWorldInformation := queryDBForSpecificWorldInformation(newWorldInformation, db)
	if entryExistsInDB {
		// If it is, hits++ and update db
		incrementHitsOnExistingWorld(&existingWorldInformation, db)
		if err = json.NewEncoder(w).Encode(existingWorldInformation); err != nil {
			err = createAndLogCustomError(err, "encode existing world info")
		}

		return
	}

	// Else add it to the db
	newWorldInformation.Hits = 1
	addNewWorldInformation(newWorldInformation, db)
	if err = json.NewEncoder(w).Encode(newWorldInformation); err != nil {
		err = createAndLogCustomError(err, "encode new world info")
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
			clearDBOnServerReset(db)
			lastResetTimeUnix = liveLastResetTimeUnix
		}
	}
}

func initLastResetTimeUnix() {
	getJS5ServerInfo(&lastResetTimeUnix)
}
