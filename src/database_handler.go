package main

import (
	"database/sql"
)

// ------------------------------------------------------------------------- //
// ----------------------------- Database Init ----------------------------- //
// ------------------------------------------------------------------------- //

// Init Database
func initDatabase() {
	// Database
	db, _ = sql.Open("sqlite3", "./ToGWorldInformation.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS World_Information (wi_id INTEGER PRIMARY KEY AUTOINCREMENT, world_number INTEGER, hits INTEGER, stream_order TEXT)")
	statement.Exec()

	statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS IP_World_Blacklist (ip_world_hash INTEGER PRIMARY KEY)")
	statement.Exec()
}

// ------------------------------------------------------------------------- //
// ----------------------------- World Info DB ----------------------------- //
// ------------------------------------------------------------------------- //

// Get information for ALL worlds
func queryDBForAllWorldInformation(database *sql.DB) []WorldInformation {

	statement, _ := database.Prepare("SELECT world_number, hits, stream_order FROM World_Information")
	rows, err := statement.Query()
	if err != nil {
		err = createAndLogCustomError(err, "Error retrieving from db in queryDBForAllWorldInformation.")
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
		err = createAndLogCustomError(err, "Error retrieving from db in queryDBForSpecificWorldInformation.")
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
		err = createAndLogCustomError(err, "There were more than one rows that had the same world_number + stream_order combination. This should not have happened.")
	}

	entryExistsInDB := numberOfEntries == 1

	return entryExistsInDB, existingWorldInformation
}

// (Update) Increment hits on existing world
func incrementHitsOnExistingWorld(worldInformation *WorldInformation, database *sql.DB) {
	worldInformation.Hits++
	statement, _ := database.Prepare("UPDATE World_Information SET hits=(?) WHERE world_number=(?) AND stream_order=(?)")
	result, err := statement.Exec(worldInformation.Hits, worldInformation.WorldNumber, worldInformation.StreamOrder)

	if err != nil {
		err = createAndLogCustomError(err, "Error in incrementHitsOnExistingWorld.")
	}

	if numRowsAffected, _ := result.RowsAffected(); numRowsAffected != 1 {
		err = createAndLogCustomError(err, "More than one rows (or 0) were updated in incrementHitsOnExistingWorld.")
	}

}

// Create new world
func addNewWorldInformation(worldInformation WorldInformation, database *sql.DB) {
	statement, _ := database.Prepare("INSERT INTO World_Information (world_number, hits, stream_order) VALUES ((?), (?), (?))")
	_, err := statement.Exec(worldInformation.WorldNumber, 1, worldInformation.StreamOrder)
	if err != nil {
		err = createAndLogCustomError(err, "Error adding to db in addNewWorldInformation.")
	}

}

// ------------------------------------------------------------------------- //
// ---------------------- IP_World_Blacklist Handling ---------------------- //
// ------------------------------------------------------------------------- //

func hasIPAlreadySubmittedDataForWorld(ipAddress string, worldNumber int, database *sql.DB) (bool, error) {
	ipWorldHash := hashIPAndWorldInfo(ipAddress, worldNumber)
	loopCount := 0
	var err error
	database.QueryRow("SELECT COUNT(ip_world_hash) FROM IP_World_Blacklist WHERE ip_world_hash=(?)", ipWorldHash).Scan(&loopCount)

	if loopCount == 0 {
		return false, err

	} else if loopCount == 1 {
		// This is the valid case where the IP has already one submission, and we should return no error and true (already submitted)

	} else if loopCount > 1 {
		err = createAndLogCustomError(err, "IP Address has more than one submissions...")

	} else {
		err = createAndLogCustomError(err, "This should be impossible. Code 1")
	}

	return true, err
}

func addIPAndWorldToDB(worldInformation WorldInformation, ipAddress string, database *sql.DB) {
	ipWorldHash := hashIPAndWorldInfo(ipAddress, worldInformation.WorldNumber)
	statement, _ := database.Prepare("INSERT INTO IP_World_Blacklist (ip_world_hash) Values ((?))")
	_, err := statement.Exec(ipWorldHash)
	if err != nil {
		err = createAndLogCustomError(err, "Error inserting into db in addIPAndWorldtoDB.")
	}
}
