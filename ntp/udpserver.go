package ntp

import (
	"fmt"
	"log/slog"
	"net"
)

type UdpServer struct {
	addr     string
	conn     *net.UDPConn
	recvChan chan *Packet
	sendChan chan *Packet
}

func NewUdpServer(port uint16) *UdpServer {
	return &UdpServer{
		addr:     fmt.Sprintf(":%d", port),
		recvChan: make(chan *Packet, 10),
		sendChan: make(chan *Packet, 10),
	}
}

func (s *UdpServer) Open() error {
	defer s.Close()
	// Listen
	addr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return err
	}
	s.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	// Write
	go func() {
		for p := range s.sendChan {
			_, err = s.conn.WriteTo(p.Data, p.Addr)
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}()
	// Read
	for {
		data := make([]byte, DataSize)
		n, addr, err := s.conn.ReadFrom(data)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		s.recvChan <- &Packet{
			Addr: addr,
			Data: data[:n],
		}
	}
}

func (s *UdpServer) Close() error {
	if s.conn == nil {
		return nil
	}
	return s.conn.Close()
}

func (s *UdpServer) Recv() <-chan *Packet {
	return s.recvChan
}

func (s *UdpServer) Send(packet *Packet) {
	s.sendChan <- packet
}
