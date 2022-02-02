package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
)

func main() {

	addr := "oldschool2.runescape.com:43594"
	//addr := "webcode.me:80"

	u := url.URL{Scheme: "ws", Host: addr}
	log.Printf("connecting to %s", u.String())

	//c, r, e := websocket.DefaultDialer.Dial(u.String(), nil)
	//fmt.Println(r)
	c, e := net.Dial("tcp", addr)
	defer c.Close()

	//_, e = c.Write([]byte("15"))
	//write24(c)

	req := "HEAD / HTTP/1.0\r\n\r\n"
	//req := ""
	_, e = c.Write([]byte(req))
	resp, e := ioutil.ReadAll(c)

	//resp, e := c.Read([]byte("255"))

	fmt.Println(c)
	fmt.Println(string(resp))
	fmt.Println(e)

}

func write24(c net.Conn) {
	pid := 255<<16 | 255
	shift16 := pid >> 16
	shift8 := pid >> 8

	c.Write([]byte(string(shift16)))
	c.Write([]byte(string(shift8)))
	c.Write([]byte(string(pid)))
}
