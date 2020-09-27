package io

import (
	"bytes"
	"encoding/json"
	"net"
	"time"

	"github.com/egoon/hanabi-server/pkg/model"
	log "github.com/sirupsen/logrus"
)

type GameStateReader interface {
	ReadGameState() (*model.GameState, error)
	Close() error
}

type gameStateReader struct {
	conn    net.Conn
	buffer  []byte
	lastIdx int
}

func NewGameStateReader(conn net.Conn) GameStateReader {
	return &gameStateReader{
		conn:   conn,
		buffer: make([]byte, 200),
	}
}

func (r *gameStateReader) ReadGameState() (*model.GameState, error) {
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
		return r.ReadGameState()
	}
	defer r.rewindBuffer(idx)
	r.lastIdx += bytesRead - idx
	GameState := model.GameState{}
	err := json.Unmarshal(r.buffer[:idx], &GameState)
	if err != nil {
		return nil, err
	}
	return &GameState, nil
}

func (r *gameStateReader) rewindBuffer(idx int) {
	i := 0
	for ; i < len(r.buffer)-idx; i++ {
		r.buffer[i] = r.buffer[i+idx]
	}
	for ; i < len(r.buffer); i++ {
		r.buffer[i] = 0
	}
}

func (r *gameStateReader) Close() error {
	return r.conn.Close()
}
