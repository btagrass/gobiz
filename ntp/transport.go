package ntp

import (
	"net"
)

var (
	DataSize = 65535 - 20 - 8 // Data size = IPv4 max size - IPv4 header size - UDP header size
)

type Packet struct {
	Addr net.Addr
	Data []byte
}

type ITransport interface {
	Open() error
	Close() error
	Recv() <-chan *Packet
	Send(packet *Packet)
}

type IClient interface {
	ITransport
	Keepalive(packet *Packet)
}

type IServer interface {
	ITransport
}

func NewClient(network, host string, port uint16) IServer {
	if network == "tcp" {
		return NewTcpClient(host, port)
	} else if network == "udp" {
		return NewUdpClient(host, port)
	}
	return nil
}

func NewServer(network string, port uint16) IServer {
	if network == "tcp" {
		return NewTcpServer(port)
	} else if network == "udp" {
		return NewUdpServer(port)
	}
	return nil
}
