package gochess

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
)

type MessageType int

const (
	Error MessageType = iota
	Register
	Deregister
	Ready
	Command
	Result
	Board
)

type WorkerStatus int

const (
	Idle WorkerStatus = iota
	Busy
)

type Message struct {
	Type MessageType `json:"type"`
	Body string      `json:"body"`
}

type Solution struct {
	Task string  `json: "Task""`
	Eval float64 `json:"Eval"`
}

type Worker struct {
	Connection *websocket.Conn
	Status     WorkerStatus
}

func ReadMessage(conn *websocket.Conn) (Message, error) {
	_, receivedBytes, err := conn.ReadMessage()
	if err != nil {
		return Message{Error, ""}, errors.New("Failed to read message.")
	}
	var receivedMessage Message
	err = json.Unmarshal(receivedBytes, &receivedMessage)
	if err != nil {
		return Message{Error, ""}, errors.New("Failed to unmarshal message.")
	}
	return receivedMessage, nil
}
