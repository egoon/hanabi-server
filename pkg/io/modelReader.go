package io

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/egoon/hanabi-server/pkg/model"
	log "github.com/sirupsen/logrus"
)

type GameStateReader interface {
	ReadGameState() (*model.GameState, error)
	Close() error
}

type ActionReader interface {
	ReadAction() (*model.Action, error)
	Close() error
}

type modelReader struct {
	conn    net.Conn
	buffer  []byte
	lastIdx int
}

func NewGameStateReader(conn net.Conn) GameStateReader {
	return &modelReader{
		conn:   conn,
		buffer: make([]byte, 1000),
	}
}

func NewActionReader(conn net.Conn) ActionReader {
	return &modelReader{
		conn:   conn,
		buffer: make([]byte, 200),
	}
}

const readTimeout = 30 //seconds

func (r *modelReader) read(model interface{}) error {
	bytesRead := 0
	if r.lastIdx == 0 || bytes.IndexByte(r.buffer, '\n') < 0 {
		err := r.conn.SetReadDeadline(time.Now().Add(time.Second * readTimeout))
		if err != nil {
			log.Warn("set read deadline failed")
		}
		bytesRead, err = r.conn.Read(r.buffer[r.lastIdx:])
		if err != nil {
			return err
		}
	}
	idx := bytes.IndexByte(r.buffer, '\n')
	if idx == -1 {
		r.lastIdx += bytesRead
		return r.read(model)
	}
	defer r.rewindBuffer(idx + 1)
	r.lastIdx += bytesRead - idx - 1
	err := json.Unmarshal(r.buffer[:idx], model)
	if err != nil {
		return fmt.Errorf("unmarshalling failed: %w", err)
	}
	return nil
}

func (r *modelReader) rewindBuffer(idx int) {
	i := 0
	for ; i < len(r.buffer)-idx; i++ {
		r.buffer[i] = r.buffer[i+idx]
	}
	for ; i < len(r.buffer); i++ {
		r.buffer[i] = 0
	}
}

func (r *modelReader) ReadGameState() (*model.GameState, error) {
	state := model.GameState{}
	err := r.read(&state)
	if err != nil {
		return nil, fmt.Errorf("failed to read game state: %w", err)
	}
	return &state, nil
}

func (r *modelReader) ReadAction() (*model.Action, error) {
	action := model.Action{}
	err := r.read(&action)
	if err != nil {
		return nil, fmt.Errorf("failed to read action: %w", err)
	}
	return &action, nil
}

func (r *modelReader) Close() error {
	return r.conn.Close()
}
