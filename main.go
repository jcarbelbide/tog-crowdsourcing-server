package main

import (
	"database/sql"
	"encoding/json"
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
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS World_Information (id INTEGER PRIMARY KEY, world_number INTEGER, hits INTEGER, stream_order TEXT)")
	statement.Exec()

	// Init Router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	//r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/worldInformation", getWorldInformation).Methods("GET")
	r.HandleFunc("/worldInformation", postWorldInformation).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getWorldInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	statement, _ := db.Prepare("SELECT world_number, hits, stream_order FROM World_Information")
	rows, err := statement.Query()
	if err != nil {
		handleError("Error retrieving from db.")
	}
	defer rows.Close()

	var worldInformationList []WorldInformation

	for rows.Next() {
		worldInformation := WorldInformation{}
		rows.Scan(&worldInformation.WorldNumber, &worldInformation.Hits, &worldInformation.StreamOrder)
		worldInformationList = append(worldInformationList, worldInformation)
	}

	//fmt.Println(worldInformationList)

	json.NewEncoder(w).Encode(worldInformationList)
	return
}

func postWorldInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var worldInformation WorldInformation
	json.NewDecoder(r.Body).Decode(&worldInformation)

	dataIsValid := verifyData(worldInformation.StreamOrder)

	fmt.Println(dataIsValid)

}

func handleError(message string) {

}

func verifyData(str string) bool {
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
