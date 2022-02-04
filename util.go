package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// -------------------------------------------------------------------------- //
// ---------------------------- Helper Functions ---------------------------- //
// -------------------------------------------------------------------------- //

func createAndLogCustomError(err error, message string) error {
	newErr := fmt.Errorf(message, err)
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

func hashIPAddress(ip string, salt string) uint64 {
	// Because we do not have a public key (username) to search the database with, if we store ips with random salts,
	// it would require us to go through every single entry in the db to find a match, since rehashing ips wont result in
	// the same hash. Therefore, we will just use a private (secret) salt value that will be left on the server. I know
	// this is about as good as not using a salt at all in the context of IP addresses, but its the best we can do.

	// I know this is absolutely not secure but it may at least slow down trolls by a few hours.
	ipByteSalted := []byte(ip + salt)

	h := sha256.New()
	h.Write(ipByteSalted)
	hashedIP := h.Sum(nil)
	intHashedIP := binary.BigEndian.Uint64(hashedIP)

	return intHashedIP
}

func hashIPAndWorldInfo(ip string, worldNumber int) uint64 {
	// First salt with private salt
	privateSaltIP := hashIPAddress(ip, privateSalt)
	// Then salt with world number. Now ip+worldNumber should be unique and can be used as key
	worldNumberSaltIP := hashIPAddress(strconv.FormatUint(privateSaltIP, 10), strconv.Itoa(worldNumber))
	return worldNumberSaltIP
}
