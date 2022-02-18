package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// -------------------------------------------------------------------------- //
// ---------------------------- Helper Functions ---------------------------- //
// -------------------------------------------------------------------------- //

func createAndLogCustomError(err error, message string) error {
	newErr := fmt.Errorf(message+" %w", err)
	log.Println(newErr)
	return newErr
}

func initLogging() *os.File {
	file, err := os.OpenFile("./info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	return file
}

func getRemoteIPAddressFromRequest(r *http.Request) string {
	return r.Header.Get("X-FORWARDED-FOR") // It is this because our NGinX server is forwarding the remote IP to us.
}

// -------------------------------------------------------------------------- //
// ---------------------------- Hashing IP Info ----------------------------- //
// -------------------------------------------------------------------------- //

func hashIPAddress(ip string, salt string) string {
	// Because we do not have a public key (username) to search the database with, if we store ips with random salts,
	// it would require us to go through every single entry in the db to find a match, since rehashing ips wont result in
	// the same hash. Therefore, we will just use a private (secret) salt value that will be left on the server. I know
	// this is about as good as not using a salt at all in the context of IP addresses, but its the best we can do.

	// I know this is absolutely not secure but it may at least slow down trolls by a few hours.
	ipByteSalted := []byte(ip + salt)

	h := sha256.New()
	h.Write(ipByteSalted)
	hashedIP := h.Sum(nil)
	//intHashedIP := binary.BigEndian.Uint64(hashedIP)			// This gets a hash in number form
	stringHashedIP := hex.EncodeToString(hashedIP)

	return stringHashedIP
}

func hashIPAndWorldInfo(ip string, worldNumber int) string {
	// First salt with private salt
	privateSaltIP := hashIPAddress(ip, privateSalt)
	// Then salt with world number. Now ip+worldNumber should be unique and can be used as key
	worldNumberSaltIP := hashIPAddress(privateSaltIP, strconv.Itoa(worldNumber))
	return worldNumberSaltIP
}

// -------------------------------------------------------------------------- //
// --------------------------- Get JS5 ServerInfo ---------------------------- //
// -------------------------------------------------------------------------- //
func getJS5ServerInfo(lastResetTime *int64) {
	resp, err := http.Get("http://localhost:8081/lastreset")

	if err != nil {
		createAndLogCustomError(err, "Error getting JS5 Server Info.")
		return
	}

	if resp == nil {
		createAndLogCustomError(err, "Resp was nil in getJS5ServerInfo.")
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		createAndLogCustomError(err, "Error reading response body when getting JS5 ServerInfo.")
		return
	}

	lastResetTimeUnixFloat, err := getJQObject(body).Float("last_reset_time_unix")
	if err != nil {
		createAndLogCustomError(err, "Error grabbing string Unix time from JSON.")
		return
	}

	*lastResetTime = int64(lastResetTimeUnixFloat)
}

func getJQObject(body []byte) *jsonq.JsonQuery {
	jsonString := string(body)
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(jsonString))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	return jq
}
