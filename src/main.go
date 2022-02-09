package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"tog-crowdsourcing-server/js5client"
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

	// Start JS5 Monitor
	go js5client.MonitorJS5Server()

	// Database
	initDatabase()

	// Init Router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	r.HandleFunc("/worldinformation", getWorldInformation).Methods("GET")
	r.HandleFunc("/worldinformation", postWorldInformation).Methods("POST")

	r.HandleFunc("/lastreset", getLastReset).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

// -------------------------------------------------------------------------- //
// ---------------------------- Request Handlers ---------------------------- //
// -------------------------------------------------------------------------- //

func getLastReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	lastReset := js5client.LastServerResetInfo
	json.NewEncoder(w).Encode(lastReset)
	return
}

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
