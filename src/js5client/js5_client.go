package js5client

import (
	"database/sql"
	"fmt"
	"time"
	"tog-crowdsourcing-server/js5connection"
)

type ServerResetInfo struct {
	LastResetTime         time.Time `json:"reset_time"`
	SecondsSinceLastReset int64     `json:"seconds_since_last_reset"`
}

var LastServerResetInfo *ServerResetInfo

func consoleLoop() {

	js5, err := js5connection.New()
	fmt.Println(err)
	//fmt.Println(js5.Ping())

	var loopCounter float64 = 0
	for {
		resp, err := js5.Ping()
		elapsedTime := js5connection.PingInterval.Seconds() * loopCounter
		fmt.Println("Resp ", resp[:10], " |  Error: ", err, " |  Seconds Elapsed: ", elapsedTime)
		if err != nil {
			fmt.Println("Connection broken! at " + time.Now().String())
			break
		}
		time.Sleep(js5connection.PingInterval)
		loopCounter++
	}

}

// MonitorJS5Server
//  1. Start by initializing db and getting the first server reset. If it's the first time, just set time to now.
//	2. Initialize JS5 Connection. If it fails to connect, keep trying. Server may be down.
//		3. Ping the connection every 5 seconds in a new loop.
//		4. If the connection breaks, break out of current loop
//	5. Be sure to take the difference between the time for last reset and current reset.
//	6. Write to DB
//	7. Set LastServerResetInfo to current one.
//	8. restart loop at step 2.
func MonitorJS5Server() {

	var db *sql.DB = initDatabase()

	LastServerResetInfo = initServerResetInfo(db)

	for { // Infinite Loop
		time.Sleep(js5connection.PingInterval)
		js5, err := js5connection.New()
		fmt.Println(&js5)

		if err != nil {
			// Start over by trying a new connection
			continue
		}

		for { // Pinging Loop. Ping until connection drops
			_, err := js5.Ping()
			if err != nil {
				break
			}
			//var input string
			//fmt.Scanln(&input)
			//if input == "y" {
			//	break
			//}
			time.Sleep(js5connection.PingInterval)
		}

		currentServerResetInfo := ServerResetInfo{
			LastResetTime:         time.Now(),
			SecondsSinceLastReset: time.Now().Unix() - LastServerResetInfo.LastResetTime.Unix(),
		}

		addNewServerResetInfo(currentServerResetInfo, db)

		LastServerResetInfo = &currentServerResetInfo
	}
}

func initServerResetInfo(database *sql.DB) *ServerResetInfo {

	var lastServerReset ServerResetInfo

	entryExists, lastServerReset := queryDBForLastServerReset(database)

	if !entryExists {
		lastServerReset = ServerResetInfo{
			LastResetTime:         time.Now(),
			SecondsSinceLastReset: 0,
		}
	}

	return &lastServerReset

}
