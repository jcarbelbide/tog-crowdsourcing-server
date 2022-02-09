package main

import (
	"database/sql"
	"net"
)

// -------------------------------------------------------------------------- //
// ------------------------- Verification Functions ------------------------- //
// -------------------------------------------------------------------------- //

func verifyDataIsValid(worldInformation WorldInformation) bool {
	worldNumberIsValid := worldInformation.WorldNumber > MIN_WORLD_VALUE && worldInformation.WorldNumber < MAX_WORLD_VALUE
	streamOrderIsValid := verifyStreamOrderIsValid(worldInformation.StreamOrder)

	return worldNumberIsValid && streamOrderIsValid
}

func verifyStreamOrderIsValid(str string) bool {
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

func verifyIPAddressIsValid(worldInformation WorldInformation, remoteIPAddress string, database *sql.DB) (bool, error) {

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
	// Things to watch out for: if IP is "" or some other garbage, should it be added to list? - No, only add valid IPs

	var err error

	// First check if the ip address is a valid IPv4 address
	if net.ParseIP(remoteIPAddress) == nil {
		err = createAndLogCustomError(nil, "Not valid IPv4 address")
		return false, err
	}

	// Check if IP has already submitted an order for that world. Do not care about the order, just care about the world. Each IP gets 1 submission per world.

	// If IP + world combo exists, return false. IP address is not valid.
	ipAlreadySubmittedDataForWorld, err := hasIPAlreadySubmittedDataForWorld(remoteIPAddress, worldInformation.WorldNumber, db)
	if ipAlreadySubmittedDataForWorld {
		return false, err

	} else {
		// Else, the IP has not submitted data for this world, so let it.
		return true, err
	}
}
