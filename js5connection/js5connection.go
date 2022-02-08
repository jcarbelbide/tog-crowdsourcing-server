package js5connection

import (
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
	WriteByte() error

	WriteInt() error

	testInitialConnection() (int, error)

	setWriteTimeout(duration time.Duration)

	setReadTimeout(duration time.Duration)

	Ping() error
}

type js5conn struct {
	conn         net.Conn
	PingInterval time.Duration
	timeout      time.Duration
}

func (c *js5conn) Ping() error {
	err := c.writeByte(255)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(c.conn)

	return err

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

func (c *js5conn) TestInitialConnection() ([]byte, error) {
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
