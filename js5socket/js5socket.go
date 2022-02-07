package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

func main() {

	addr := "oldschool2.runescape.com:43594"
	timeoutDuration := 5000 * time.Millisecond
	sleepDuration := 5000 * time.Millisecond

	fmt.Println(addr)
	conn, err := net.Dial("tcp", addr)
	defer conn.Close()

	if err != nil {
		fmt.Println("Error establishing connection")
	}

	var req byte = 255
	var reqArray []byte

	for {
		reqArray = append(reqArray, req)
		conn.SetWriteDeadline(time.Now().Add(timeoutDuration))
		_, err = conn.Write(reqArray)

		conn.SetReadDeadline(time.Now().Add(timeoutDuration))
		resp, err := ioutil.ReadAll(conn)

		fmt.Println(conn)
		fmt.Println(string(resp))
		fmt.Println(err)

		if err != nil {
			fmt.Println("Connection broken!")
			break
		}

		time.Sleep(sleepDuration)

	}

}

func write24(c net.Conn) {
	pid := 255<<16 | 255
	shift16 := pid >> 16
	shift8 := pid >> 8

	c.Write([]byte(string(shift16)))
	c.Write([]byte(string(shift8)))
	c.Write([]byte(string(pid)))
}
