package main

import (
	"database/sql"
)

// ------------------------------------------------------------------------- //
// --------------------------- Database Handling --------------------------- //
// ------------------------------------------------------------------------- //

// Get information for ALL worlds
func queryDBForAllWorldInformation(database *sql.DB) []WorldInformation {

	statement, _ := database.Prepare("SELECT world_number, hits, stream_order FROM World_Information")
	rows, err := statement.Query()
	if err != nil {
		err = createCustomError(err, "Error retrieving from db in queryDBForAllWorldInformation.")
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
		err = createCustomError(err, "Error retrieving from db in queryDBForSpecificWorldInformation.")
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
		err = createCustomError(err, "There were more than one rows that had the same world_number + stream_order combination. This should not have happened.")
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
		err = createCustomError(err, "Error in incrementHitsOnExistingWorld.")
	}

	if numRowsAffected, _ := result.RowsAffected(); numRowsAffected != 1 {
		err = createCustomError(err, "More than one rows (or 0) were updated in incrementHitsOnExistingWorld.")
	}

}

// Create new world
func addNewWorldInformation(worldInformation WorldInformation, database *sql.DB) {
	statement, _ := db.Prepare("INSERT INTO World_Information (world_number, hits, stream_order) VALUES ((?), (?), (?))")
	_, err := statement.Exec(worldInformation.WorldNumber, 1, worldInformation.StreamOrder)
	if err != nil {
		err = createCustomError(err, "Error retrieving from db in queryDBForSpecificWorldInformation.")
	}

}

func hasIPAlreadySubmittedDataForWorld(worldNumber int, ipAddress string, database *sql.DB) (bool, error) {
	statement, _ := database.Prepare("SELECT world_number, ip_address FROM IP_List WHERE world_number=(?) AND ip_address=(?)")
	rows, err := statement.Query(worldNumber, ipAddress)

	if err != nil {
		err = createCustomError(err, "Error retrieving from db in hasIPAlreadySubmittedDataForWorld.")
	}
	defer rows.Close()

	var dbWorldNumber int
	var dbIPAddress string
	loopCount := 0

	for rows.Next() {
		rows.Scan(&dbWorldNumber, &dbIPAddress)
		loopCount++
	}

	if loopCount == 0 {
		return false, err

	} else if loopCount > 1 {
		err = createCustomError(err, "IP Address has more than one submissions...")

	} else if loopCount == 1 {
		// This is the valid case where the IP has already one submission and we should return no error and true (already submitted)

	} else {
		err = createCustomError(err, "This should be impossible. Code 1")
	}

	return true, err
}
