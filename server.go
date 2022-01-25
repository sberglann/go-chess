package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strings"
)

type MoveResponse struct {
	Move       string   `json:"move"`
	LegalMoves []string `json:"legalMoves"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

func StartServer() {
	http.HandleFunc("/chess", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		for {
			// Read message from browser
			msgType, receivedMessage, err := conn.ReadMessage()
			if err != nil {
				return
			}
			receivedMessageString := string(receivedMessage)

			var legalFenMoves []string
			var nextMove BitBoard
			if receivedMessageString == "init" {
				nextMove = StartBoard
			} else {
				board := BoardFromFEN(receivedMessageString)
				legalMoves := GenerateLegalMoves(board)
				// For now, just pick a random move.
				nextMove = legalMoves[rand.Intn(len(legalMoves))]
			}

			legalMoves := GenerateLegalMoves(nextMove)

			for _, move := range legalMoves {
				onlyPieces := strings.Split(move.ToFEN(), " ")[0]
				legalFenMoves = append(legalFenMoves, onlyPieces)
			}
			response := &MoveResponse{nextMove.ToFEN(), legalFenMoves}
			message, err := json.Marshal(response)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(receivedMessage))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, message); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "gui/index.html")
	})

	http.ListenAndServe(":8080", nil)
}
