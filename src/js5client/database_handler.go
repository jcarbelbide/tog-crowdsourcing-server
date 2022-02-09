package js5client

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

// ------------------------------------------------------------------------- //
// ----------------------------- Database Init ----------------------------- //
// ------------------------------------------------------------------------- //

// Init Database
func initDatabase() *sql.DB {
	// Database
	database, _ := sql.Open("sqlite3", "./ServerResetTimes.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Server_Reset_Times (row_id INTEGER PRIMARY KEY AUTOINCREMENT, last_reset_time TEXT, seconds_since_last_reset INTEGER)")
	statement.Exec()
	return database
}

// ------------------------------------------------------------------------- //
// ----------------------------- World Info DB ----------------------------- //
// ------------------------------------------------------------------------- //

// Get information for specific world
func queryDBForLastServerReset(database *sql.DB) (bool, ServerResetInfo) {

	statement, _ := database.Prepare("SELECT last_reset_time, seconds_since_last_reset FROM Server_Reset_Times ORDER BY row_id DESC LIMIT 1")
	rows, err := statement.Query()

	if err != nil {
		err = createAndLogCustomError(err, "Error retrieving from db in queryDBForLastServerReset.")
	}
	defer rows.Close()

	numberOfEntries := 0
	var serverResetInfo ServerResetInfo
	var intSecondsSinceLastReset int
	var textLastReset string
	for rows.Next() {
		rows.Scan(&textLastReset, &intSecondsSinceLastReset)
		serverResetInfo = ServerResetInfo{
			LastResetTime:         time.Unix(time.Now().Unix()-int64(intSecondsSinceLastReset), 0), // TODO This is technically wrong. We should probably just save the Unix time along with the text to create the last server reset.
			SecondsSinceLastReset: int64(intSecondsSinceLastReset),
		}
		numberOfEntries++
	}

	entryExistsInDB := numberOfEntries == 1

	return entryExistsInDB, serverResetInfo
}

// Create new world
func addNewServerResetInfo(serverResetInfo ServerResetInfo, database *sql.DB) {
	statement, _ := database.Prepare("INSERT INTO Server_Reset_Times (last_reset_time, seconds_since_last_reset) VALUES ((?), (?))")
	_, err := statement.Exec(serverResetInfo.LastResetTime.String(), serverResetInfo.SecondsSinceLastReset)
	if err != nil {
		err = createAndLogCustomError(err, "Error adding to db in addNewServerResetInfo.")
	}

}

// ------------------------------------------------------------------------- //
// --------------------------- Helper Functions ---------------------------- //
// ------------------------------------------------------------------------- //

func createAndLogCustomError(err error, message string) error {
	newErr := fmt.Errorf(message+" %w", err)
	log.Println(newErr)
	return newErr
}
