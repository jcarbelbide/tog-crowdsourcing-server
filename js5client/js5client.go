package main

import (
	"fmt"
	"net"
	"time"
	"tog-crowdsourcing-server/js5connection"
)

func main() {

	js5 := js5connection.New()
	fmt.Println(js5.TestInitialConnection())
	var loopCounter float64 = 0
	for {
		err := js5.Ping()
		elapsedTime := js5.PingInterval.Seconds() * loopCounter
		fmt.Println("Error: ", err, " |  Seconds Elapsed: ", elapsedTime)
		if err != nil {
			fmt.Println("Connection broken! at " + time.Now().String())
			break
		}
		time.Sleep(js5.PingInterval)
		loopCounter++
	}

	//sleepDuration := 5000 * time.Millisecond
	//addr := "oldschool2.runescape.com:43594"
	//
	//fmt.Println(addr)
	//conn, err := net.Dial("tcp", addr)
	//defer conn.Close()
	//
	//if err != nil {
	//	fmt.Println("Error establishing connection")
	//}
	//
	////var req byte = 0xFF
	////var rev byte = 0xff
	////var reqArray []byte
	//var resp []byte
	//
	//for {
	//	conn.SetWriteDeadline(time.Now().Add(timeoutDuration))
	//	_, err = conn.Write(intToByteArray(15))
	//
	//	conn.SetWriteDeadline(time.Now().Add(timeoutDuration))
	//	err = writeInt(0, conn)
	//
	//	conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	//	resp, err = ioutil.ReadAll(conn)
	//
	//	fmt.Println(conn)
	//	fmt.Println(resp)
	//	fmt.Println(err)
	//
	//	if err != nil {
	//		fmt.Println("Connection broken!")
	//		break
	//	}
	//
	//	time.Sleep(sleepDuration)
	//
	//}

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
