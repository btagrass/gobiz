package ntp

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
)

type TcpServer struct {
	addr     string
	listener net.Listener
	conns    sync.Map
	recvChan chan *Packet
	sendChan chan *Packet
}

func NewTcpServer(port uint16) *TcpServer {
	return &TcpServer{
		addr:     fmt.Sprintf(":%d", port),
		recvChan: make(chan *Packet, 10),
		sendChan: make(chan *Packet, 10),
	}
}

func (s *TcpServer) Open() error {
	defer s.Close()
	// Listen
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer s.listener.Close()
	// Write
	go func() {
		for p := range s.sendChan {
			v, ok := s.conns.Load(p.Addr.String())
			if !ok {
				continue
			}
			conn := v.(net.Conn)
			_, err = conn.Write(p.Data)
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}()
	// Read
	for {
		data := make([]byte, DataSize)
		conn, err := s.listener.Accept()
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		s.conns.Store(conn.RemoteAddr().String(), conn)
		n, err := conn.Read(data)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		s.recvChan <- &Packet{
			Addr: conn.RemoteAddr(),
			Data: data[:n],
		}
	}
}

func (s *TcpServer) Close() error {
	s.conns.Range(func(k, v any) bool {
		c := v.(net.Conn)
		err := c.Close()
		if err != nil {
			slog.Error(err.Error())
		}
		s.conns.Delete(k)
		return true
	})
	return nil
}

func (s *TcpServer) Recv() <-chan *Packet {
	return s.recvChan
}

func (s *TcpServer) Send(packet *Packet) {
	s.sendChan <- packet
}
