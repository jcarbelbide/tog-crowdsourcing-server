package js5client

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"time"
)

// ------------------------------------------------------------------------- //
// ----------------------------- Database Init ----------------------------- //
// ------------------------------------------------------------------------- //

// Init Database
func initDatabase() *sql.DB {
	// Database
	database, _ := sql.Open("sqlite3", "./ServerResetTimes.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS Server_Reset_Times (row_id INTEGER PRIMARY KEY AUTOINCREMENT, last_reset_time TEXT, last_reset_time_unix TEXT, last_server_uptime TEXT)")
	statement.Exec()
	return database
}

// ------------------------------------------------------------------------- //
// ----------------------------- World Info DB ----------------------------- //
// ------------------------------------------------------------------------- //

// Get information for specific world
func queryDBForLastServerReset(database *sql.DB) (bool, ServerResetInfo) {

	statement, _ := database.Prepare("SELECT last_reset_time, last_reset_time_unix, last_server_uptime FROM Server_Reset_Times ORDER BY row_id DESC LIMIT 1")
	rows, err := statement.Query()

	if err != nil {
		err = createAndLogCustomError(err, "Error retrieving from db in queryDBForLastServerReset.")
	}
	defer rows.Close()

	numberOfEntries := 0
	var serverResetInfo ServerResetInfo

	var stringLastResetTime string
	var stringLastResetTimeUnix string
	var stringLastServerUptime string

	for rows.Next() {
		rows.Scan(&stringLastResetTime, &stringLastResetTimeUnix, &stringLastServerUptime)

		int64LastResetTimeUnix, _ := strconv.ParseInt(stringLastResetTimeUnix, 10, 64)
		int64LastServerUptime, _ := strconv.ParseInt(stringLastServerUptime, 10, 64)

		serverResetInfo = ServerResetInfo{
			LastResetTime:     time.Unix(int64LastResetTimeUnix, 0),
			LastResetTimeUnix: int64LastResetTimeUnix,
			LastServerUptime:  int64LastServerUptime,
		}
		numberOfEntries++
	}

	entryExistsInDB := numberOfEntries == 1

	return entryExistsInDB, serverResetInfo
}

// Create new world
func addNewServerResetInfo(serverResetInfo ServerResetInfo, database *sql.DB) {
	statement, _ := database.Prepare("INSERT INTO Server_Reset_Times (last_reset_time, last_reset_time_unix, last_server_uptime) VALUES ((?), (?), (?))")
	_, err := statement.Exec(serverResetInfo.LastResetTime.String(), strconv.FormatInt(serverResetInfo.LastResetTimeUnix, 10), strconv.FormatInt(serverResetInfo.LastServerUptime, 10))
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
