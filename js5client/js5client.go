package main

import (
	"fmt"
	"net"
	"time"
	"tog-crowdsourcing-server/js5connection"
)

func main() {

	js5 := js5connection.New()
	fmt.Println(js5.WriteJS5Header())
	fmt.Println(js5.Ping())

	var loopCounter float64 = 0
	for {
		resp, err := js5.Ping()
		elapsedTime := js5.PingInterval.Seconds() * loopCounter
		fmt.Println("Resp ", resp, " |  Error: ", err, " |  Seconds Elapsed: ", elapsedTime)
		if err != nil {
			fmt.Println("Connection broken! at " + time.Now().String())
			break
		}
		time.Sleep(js5.PingInterval)
		loopCounter++
	}

}

func intToByteArray(num int) []byte {
	return append(make([]byte, 0), byte(num))
}

func writeInt(v int, out net.Conn) error {
	out.Write(intToByteArray((v >> 24) & 0xFF))
	out.Write(intToByteArray((v >> 16) & 0xFF))
	out.Write(intToByteArray((v >> 8) & 0xFF))
	out.Write(intToByteArray((v >> 0) & 0xFF))

	return nil
}

func write24(out net.Conn) {
	pid := 255<<16 | 255
	shift16 := pid >> 16
	shift8 := pid >> 8

	out.Write(intToByteArray(shift16))
	out.Write(intToByteArray(shift8))
	out.Write(intToByteArray(pid))
}
