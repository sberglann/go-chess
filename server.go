package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type MoveResponse struct {
	Move       string              `json:"move"`
	LegalMoves []LegalMoveResponse `json:"legalMoves"`
	Eval       float64             `json:"eval"`
}

type LegalMoveResponse struct {
	ClientFen    string `json:"clientFen"`
	TrueFen      string `json:"trueFen"`
	IsCastleMove bool   `json:"IsCastleMove"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

var server *http.Server

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

			var legalStates []LegalMoveResponse
			var nextMove BitBoard
			var eval float64
			if receivedMessageString == "init" {
				nextMove = StartBoard
			} else if receivedMessageString == "quit" {
				server.Close()
			} else {
				board := BoardFromFEN(receivedMessageString)
				// No time constraint for web server - use default depth
				evaluatedBoard := BestMove(board, 0)
				nextMove = evaluatedBoard.board
				eval = evaluatedBoard.eval

			}

			legalMoves, _ := GenerateLegalStates(nextMove)

			for _, move := range legalMoves {
				lmr := extractLegalMoveResponse(nextMove, move)
				legalStates = append(legalStates, lmr)
			}
			response := &MoveResponse{nextMove.ToFEN(), legalStates, eval}
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

	server = &http.Server{
		Addr:         ":8081",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func extractLegalMoveResponse(previous BitBoard, next BitBoard) LegalMoveResponse {
	whiteKingWasInPosition := previous.KingBB&previous.WhiteBB&posToBitBoard(4) > 0
	wkCastle := next.KingBB&next.WhiteBB&posToBitBoard(6) > 0
	wqCastle := next.KingBB&next.WhiteBB&posToBitBoard(2) > 0
	isCastleMove := false
	trueFen := next.ToFEN()

	if whiteKingWasInPosition && wkCastle && previous.WhiteCanCastleKingSite() {
		next.RookBB &^= posToBitBoard(5)
		next.RookBB |= posToBitBoard(7)
		next.WhiteBB &^= posToBitBoard(5)
		next.WhiteBB |= posToBitBoard(7)
		isCastleMove = true
	} else if whiteKingWasInPosition && wqCastle && previous.WhiteCanCastleQueenSite() {
		next.RookBB &^= posToBitBoard(3)
		next.RookBB |= posToBitBoard(0)
		next.WhiteBB &^= posToBitBoard(3)
		next.WhiteBB |= posToBitBoard(0)
		isCastleMove = true
	}
	clientFen := strings.Split(next.ToFEN(), " ")[0]

	return LegalMoveResponse{TrueFen: trueFen, ClientFen: clientFen, IsCastleMove: isCastleMove}
}
