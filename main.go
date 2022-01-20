package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
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
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS world_information (id INTEGER PRIMARY KEY, worldNumber INTEGER, hits INTEGER, streamOrder TEXT)")
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

	statement, _ := db.Prepare("SELECT worldNumber, hits, streamOrder FROM world_information")
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

	fmt.Println(worldInformationList)

	json.NewEncoder(w).Encode(worldInformationList)
	return
}

func postWorldInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params. Should be world number only. Verify that. TODO: Change this so that it gets from body, not params.

	worldNumber, err := strconv.Atoi(params["world_number"])
	if err != nil {
		handleError("Error converting world number from request params to int.")
	}

	// TODO: Verify data is good (3g 3b) before adding to db.

	worldNumber = worldNumber * 2 / 2
}

func handleError(message string) {

}
