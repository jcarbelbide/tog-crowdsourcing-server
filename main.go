package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

type WorldInformation struct {
	WorldNumber int    `json:"world_number"`
	Hits        int    `json:"hits"`
	StreamOrder string `json:"stream_order"`
}

var db *sql.DB

func main() {

	// Database
	db, _ = sql.Open("sqlite3", "./ToGWorldInformation.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS World_Information (wi_id INTEGER PRIMARY KEY AUTOINCREMENT, world_number INTEGER, hits INTEGER, stream_order TEXT)")
	statement.Exec()

	statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS IP_List (ip_id INTEGER PRIMARY KEY AUTOINCREMENT, world_number INTEGER, ip_address TEXT)")
	statement.Exec()

	// Init Router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	//r.HandleFunc("/api/books", getBooks).Methods("GET")
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

	verifyIPAddressIsValid(r, db)
	dataIsValid := verifyDataIsValid(newWorldInformation.StreamOrder)

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

func handleError(err error, message string) error {
	return errors.New(message)
}

// -------------------------------------------------------------------------- //
// ------------------------- Verification Functions ------------------------- //
// -------------------------------------------------------------------------- //

func verifyDataIsValid(str string) bool {
	greenCount := 0
	blueCount := 0
	for _, c := range str {
		if c == 'g' {
			greenCount++
		} else if c == 'b' {
			blueCount++
		} else {
			return false
		}
	}
	return greenCount == 3 && blueCount == 3 && len(str) == 6
}

func verifyIPAddressIsValid(r *http.Request, database *sql.DB) bool {

	// TODO: In this function, we want to make sure that any one ip address is not sending so many requests
	// Since we cannot stop people from sending bad requests, we will have to settle for limiting the number
	// of requests one ip can send without hindering the ability for legitimate users to send data to the server.
	// Legitimate users will only have to send data to the server once per world per server reset, since the data should
	// not change until another server reset.

	// One IP address can submit one request per world per server reset.
	// In other words: 1 IP/world/reset
	// 'Per reset' will be implemented by clearing the list of IPs everytime a reset is detected.
	// 'Per world' Originally, I thought we could just keep track of the wi_id and if the ip has already sent a request
	// with that same wi_id, then don't do anything. However, then they could submit a ton of requests with different stream_orders.
	// What REALLY needs to be looked at is the world. There is really no good reason why one IP should be submitting the
	// stream_order for one world multiple times during one server reset, since it should always be the same. Thus, if an
	// IP is submitting a request for a world in which that IP has already submitted a valid request, the request should
	// be ignored.
	//
	// In order to do this, we will keep a separate table in the db with (key, world, ip)
	// When checking, we query the table where the ip and world matches the request. There should be at most only be one entry.
	// Return true if the entry DNE, and return false if the entry exists. Error if there is more than one entry.
	// Things to watch out for: if IP is "" or some other garbage, should it be added to list?

	fmt.Println("Forward: ", r.Header.Get("X-FORWARDED-FOR"))
	fmt.Println("request.RemoteAddr: ", r.RemoteAddr)

	return false
}

// ------------------------------------------------------------------------- //
// --------------------------- Database Handling --------------------------- //
// ------------------------------------------------------------------------- //

// Get information for ALL worlds
func queryDBForAllWorldInformation(database *sql.DB) []WorldInformation {

	statement, _ := database.Prepare("SELECT world_number, hits, stream_order FROM World_Information")
	rows, err := statement.Query()
	if err != nil {
		handleError(err, "Error retrieving from db in queryDBForAllWorldInformation.")
	}
	defer rows.Close()

	var worldInformationList []WorldInformation

	for rows.Next() {
		worldInformation := WorldInformation{}
		rows.Scan(&worldInformation.WorldNumber, &worldInformation.Hits, &worldInformation.StreamOrder)
		worldInformationList = append(worldInformationList, worldInformation)
	}

	return worldInformationList
}

// Get information for specific world
func queryDBForSpecificWorldInformation(worldInformation WorldInformation, database *sql.DB) (bool, WorldInformation) {

	statement, _ := database.Prepare("SELECT world_number, hits, stream_order FROM World_Information WHERE world_number=(?) AND stream_order=(?)")
	rows, err := statement.Query(worldInformation.WorldNumber, worldInformation.StreamOrder)

	if err != nil {
		handleError(err, "Error retrieving from db in queryDBForSpecificWorldInformation.")
	}
	defer rows.Close()

	numberOfEntries := 0
	var existingWorldInformation WorldInformation
	for rows.Next() {
		existingWorldInformation = WorldInformation{}
		rows.Scan(&existingWorldInformation.WorldNumber, &existingWorldInformation.Hits, &existingWorldInformation.StreamOrder)
		numberOfEntries++
	}

	if numberOfEntries > 1 {
		handleError(err, "There were more than one rows that had the same world_number + stream_order combination. This should not have happened.")
	}

	entryExistsInDB := numberOfEntries == 1

	return entryExistsInDB, existingWorldInformation
}

// (Update) Increment hits on existing world
func incrementHitsOnExistingWorld(worldInformation *WorldInformation, database *sql.DB) {
	worldInformation.Hits++
	statement, _ := db.Prepare("UPDATE World_Information SET hits=(?) WHERE world_number=(?) AND stream_order=(?)")
	result, err := statement.Exec(worldInformation.Hits, worldInformation.WorldNumber, worldInformation.StreamOrder)

	if err != nil {
		handleError(err, "Error in incrementHitsOnExistingWorld.")
	}

	if numRowsAffected, _ := result.RowsAffected(); numRowsAffected != 1 {
		handleError(err, "More than one rows (or 0) were updated in incrementHitsOnExistingWorld.")
	}

}

// Create new world
func addNewWorldInformation(worldInformation WorldInformation, database *sql.DB) {
	statement, _ := db.Prepare("INSERT INTO World_Information (world_number, hits, stream_order) VALUES ((?), (?), (?))")
	_, err := statement.Exec(worldInformation.WorldNumber, 1, worldInformation.StreamOrder)
	if err != nil {
		handleError(err, "Error retrieving from db in queryDBForSpecificWorldInformation.")
	}

}
