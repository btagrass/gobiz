package ntp

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type UdpClient struct {
	addr     string
	conn     net.Conn
	recvChan chan *Packet
	sendChan chan *Packet
}

func NewUdpClient(host string, port uint16) *UdpClient {
	return &UdpClient{
		addr:     fmt.Sprintf("%s:%d", host, port),
		recvChan: make(chan *Packet, 10),
		sendChan: make(chan *Packet, 10),
	}
}

func (c *UdpClient) Open() error {
	defer c.Close()
	// Listen
	var err error
	c.conn, err = net.Dial("udp", c.addr)
	if err != nil {
		return err
	}
	// Write
	go func() {
		for p := range c.sendChan {
			_, err = c.conn.Write(p.Data)
			if err != nil {
				logrus.Error(err)
			}
		}
	}()
	// Read
	for {
		data := make([]byte, DataSize)
		n, err := c.conn.Read(data)
		if err != nil {
			logrus.Error(err)
			continue
		}
		c.recvChan <- &Packet{
			Addr: c.conn.RemoteAddr(),
			Data: data[:n],
		}
	}
}

func (c *UdpClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *UdpClient) Recv() <-chan *Packet {
	return c.recvChan
}

func (c *UdpClient) Send(packet *Packet) {
	c.sendChan <- packet
}

func (c *UdpClient) Keepalive(packet *Packet) {
	if packet == nil {
		packet = &Packet{
			Data: []byte("ping"),
		}
	}
	c.Send(packet)
}
