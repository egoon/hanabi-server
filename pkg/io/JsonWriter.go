package io

import (
	"encoding/json"
	"net"
)

type JsonWriter interface {
	Write(interface{}) (int, error)
	Close() error
}

type jsonWriter struct {
	conn net.Conn
}

func NewJsonWriter(conn net.Conn) JsonWriter {
	return &jsonWriter{
		conn: conn,
	}
}

func (w *jsonWriter) Write(obj interface{}) (int, error) {
	msg, _ := json.Marshal(obj)
	msg = append(msg, '\n')
	return w.conn.Write(msg)
}

func (w *jsonWriter) Close() error {
	return w.conn.Close()
}
