package io

import (
	"bytes"
	"encoding/json"
	"github.com/egoon/hanabi-server/pkg/model"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type ActionReader interface {
	ReadAction() (*model.Action, error)
	Close() error
}

type actionReader struct {
	conn    net.Conn
	buffer  []byte
	lastIdx int
}

func NewActionReader(conn net.Conn) ActionReader {
	return &actionReader{
		conn:   conn,
		buffer: make([]byte, 200),
	}
}

const readTimeout = 30 //seconds

func (r *actionReader) ReadAction() (*model.Action, error) {
	bytesRead := 0
	if r.lastIdx == 0 || bytes.IndexByte(r.buffer, '\n') < 0 {
		err := r.conn.SetReadDeadline(time.Now().Add(time.Second * readTimeout))
		if err != nil {
			log.Warn("set read deadline failed")
		}
		bytesRead, err := r.conn.Read(r.buffer)
		if err != nil {
			return nil, err
		}
		log.Info("[", string(r.buffer[:bytesRead]), "]")
	}
	idx := bytes.IndexByte(r.buffer, '\n')
	if idx == -1 {
		r.lastIdx += bytesRead
		return r.ReadAction()
	}
	defer r.rewindBuffer(idx)
	r.lastIdx += bytesRead - idx
	action := model.Action{}
	err := json.Unmarshal(r.buffer[:idx], &action)
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func (r *actionReader) rewindBuffer(idx int) {
	i := 0
	for ; i < len(r.buffer)-idx; i++ {
		r.buffer[i] = r.buffer[i+idx]
	}
	for ; i < len(r.buffer); i++ {
		r.buffer[i] = 0
	}
}

func (r *actionReader) Close() error {
	return r.conn.Close()
}
