package js5connection

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

// JS5Connection
// Very rudimentary implementation of JS5 Connection in GO.
// The purpose of this JS5 Connection is simply to establish a connection
// to a single Runescape JS5 socket. The connection can then be used to
// determine when the Runescape servers have been reset. If the servers
// have been reset, then the socket connection will be closed.
type JS5Connection interface {
	WriteJS5Header() (int, error)

	Ping() error
}

type js5conn struct {
	conn         net.Conn
	PingInterval time.Duration
	timeout      time.Duration
}

func (c *js5conn) Ping() ([]byte, error) {
	err := c.writePID()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(c.conn)

}

func (c *js5conn) writeByte(data int) error {
	_, err := c.write(intToByteArray(data))
	return err
}

func (c *js5conn) writeInt(data int) error {
	_, err := c.write(intToByteArray((data >> 24) & 0xFF))
	if err != nil {
		return err
	}

	_, err = c.write(intToByteArray((data >> 16) & 0xFF))
	if err != nil {
		return err
	}

	_, err = c.write(intToByteArray((data >> 8) & 0xFF))
	if err != nil {
		return err
	}

	_, err = c.write(intToByteArray((data >> 0) & 0xFF))
	return err
}

func (c *js5conn) WriteJS5Header() ([]byte, error) {
	err := c.writeByte(15)
	if err != nil {
		return nil, err
	}

	err = c.writeInt(0)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(c.conn)

}

func (c *js5conn) setWriteTimeout(duration time.Duration) {
	c.conn.SetWriteDeadline(time.Now().Add(duration))
}

func (c *js5conn) setReadTimeout(duration time.Duration) {
	c.conn.SetReadDeadline(time.Now().Add(duration))
}

func (c *js5conn) write(b []byte) (int, error) {
	c.setWriteTimeout(c.timeout)
	return c.conn.Write(b)
}

func (c *js5conn) writePID() error {
	pid := 255<<16 | 255

	err := c.writeByte(1)
	if err != nil {
		fmt.Println("testRequest Code 1")
		return err
	}

	err = c.writeByte(pid >> 16)
	if err != nil {
		fmt.Println("testRequest Code 2")
		return err
	}

	err = c.writeByte(pid >> 8)
	if err != nil {
		fmt.Println("testRequest Code 3")
		return err
	}

	err = c.writeByte(pid)

	return err
}

func intToByteArray(num int) []byte {
	return append(make([]byte, 0), byte(num))
}

func New() *js5conn {
	addr := "oldschool2.runescape.com:43594"
	conn, _ := net.Dial("tcp", addr)

	var c = js5conn{
		conn:         conn,
		PingInterval: 5000 * time.Millisecond,
		timeout:      5000 * time.Millisecond,
	}

	return &c
}
