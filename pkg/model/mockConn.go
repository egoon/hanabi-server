package model

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type MockConn struct {
	BytesToRead        []byte
	BytesWritten       chan []byte
	DeadLine           time.Time
	ReadDeadLine       time.Time
	WriteDeadLine      time.Time
	ReadError          error
	WriteError         error
	CloseError         error
	DeadlineError      error
	ReadDeadlineError  error
	WriteDeadlineError error
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	for i := 0; i < len(m.BytesToRead); i++ {
		b[i] = m.BytesToRead[i]
	}
	return len(m.BytesToRead), m.ReadError
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	m.BytesWritten <- b
	log.Info("state written")
	return len(b), m.WriteError
}

func (m *MockConn) Close() error {
	return m.CloseError
}

func (m *MockConn) LocalAddr() net.Addr {
	return nil
}

func (m *MockConn) RemoteAddr() net.Addr {
	return nil
}

func (m *MockConn) SetDeadline(t time.Time) error {
	m.DeadLine = t
	return m.DeadlineError
}

func (m *MockConn) SetReadDeadline(t time.Time) error {
	m.ReadDeadLine = t
	return m.DeadlineError
}

func (m *MockConn) SetWriteDeadline(t time.Time) error {
	m.WriteDeadLine = t
	return m.WriteDeadlineError
}
