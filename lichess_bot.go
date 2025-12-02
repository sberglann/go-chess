package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const LichessAPIBase = "https://lichess.org"

type LichessBot struct {
	token      string
	httpClient *http.Client
	userID     string // Cache user ID to avoid repeated API calls
}

type Challenge struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Challenger struct {
		ID string `json:"id"`
	} `json:"challenger"`
	DestUser struct {
		ID string `json:"id"`
	} `json:"destUser"`
	Variant struct {
		Key string `json:"key"`
	} `json:"variant"`
	Rated bool `json:"rated"`
	TimeControl struct {
		Type string `json:"type"`
	} `json:"timeControl"`
}

type GameState struct {
	Type      string          `json:"type"`
	GameID    string          `json:"gameId,omitempty"`
	Status    json.RawMessage `json:"status,omitempty"` // Can be string or object in game stream
	Moves     string          `json:"moves,omitempty"`
	Fen       string          `json:"fen,omitempty"`
	WTime     int             `json:"wtime,omitempty"`
	BTime     int             `json:"btime,omitempty"`
	WInc      int             `json:"winc,omitempty"`
	BInc      int             `json:"binc,omitempty"`
	White     *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"white,omitempty"`
	Black *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"black,omitempty"`
	Winner string `json:"winner,omitempty"`
}

type GameFull struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Rated  bool   `json:"rated"`
	Variant struct {
		Key string `json:"key"`
	} `json:"variant"`
	Clock *struct {
		Initial   int `json:"initial"`
		Increment int `json:"increment"`
	} `json:"clock"`
	Speed string `json:"speed"`
	Perf  struct {
		Name string `json:"name"`
	} `json:"perf"`
	White struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"white"`
	Black struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"black"`
	InitialFen string `json:"initialFen"`
		State      struct {
			Type      string          `json:"type"`
			Moves     string          `json:"moves"`
			Fen       string          `json:"fen,omitempty"` // FEN might be in state for gameFull
			WTime     int             `json:"wtime"`
			BTime     int             `json:"btime"`
			WInc      int             `json:"winc"`
			BInc      int             `json:"binc"`
			Status    json.RawMessage `json:"status"` // Can be string or object
			Winner    string          `json:"winner,omitempty"`
		} `json:"state"`
}

func NewLichessBot(token string) *LichessBot {
	return &LichessBot{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (bot *LichessBot) makeRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+bot.token)
	req.Header.Set("Content-Type", "application/json")
	
	return bot.httpClient.Do(req)
}

func (bot *LichessBot) acceptChallenge(challengeID string) error {
	url := fmt.Sprintf("%s/api/challenge/%s/accept", LichessAPIBase, challengeID)
	resp, err := bot.makeRequest("POST", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to accept challenge: %s - %s", resp.Status, string(body))
	}
	
	log.Printf("Accepted challenge %s", challengeID)
	return nil
}

func (bot *LichessBot) declineChallenge(challengeID string) error {
	url := fmt.Sprintf("%s/api/challenge/%s/decline", LichessAPIBase, challengeID)
	resp, err := bot.makeRequest("POST", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to decline challenge: %s - %s", resp.Status, string(body))
	}
	
	log.Printf("Declined challenge %s", challengeID)
	return nil
}

func (bot *LichessBot) makeMove(gameID, move string) error {
	url := fmt.Sprintf("%s/api/bot/game/%s/move/%s", LichessAPIBase, gameID, move)
	resp, err := bot.makeRequest("POST", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to make move: %s - %s", resp.Status, string(body))
	}
	
	log.Printf("Made move %s in game %s", move, gameID)
	return nil
}

func (bot *LichessBot) declineDraw(gameID string) error {
	url := fmt.Sprintf("%s/api/bot/game/%s/draw/no", LichessAPIBase, gameID)
	resp, err := bot.makeRequest("POST", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to decline draw: %s - %s", resp.Status, string(body))
	}
	
	log.Printf("Declined draw offer in game %s", gameID)
	return nil
}

func (bot *LichessBot) declineTakeback(gameID string) error {
	url := fmt.Sprintf("%s/api/bot/game/%s/takeback/no", LichessAPIBase, gameID)
	resp, err := bot.makeRequest("POST", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to decline takeback: %s - %s", resp.Status, string(body))
	}
	
	log.Printf("Declined takeback offer in game %s", gameID)
	return nil
}

func (bot *LichessBot) streamEvents() error {
	url := fmt.Sprintf("%s/api/stream/event", LichessAPIBase)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "Bearer "+bot.token)
	
	// Use a client without timeout for long-running streams
	// Also disable connection pooling to avoid issues
	streamClient := &http.Client{
		Timeout: 0, // No timeout for streaming
		Transport: &http.Transport{
			DisableKeepAlives: false, // Keep connection alive for streaming
		},
	}
	
	log.Println("Opening event stream connection...")
	resp, err := streamClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to stream events: %s - %s", resp.Status, string(body))
	}
	
	log.Println("Event stream connected. Waiting for events...")
	
	// Use a larger buffer for the scanner to handle long lines
	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // 1MB buffer
	
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			if DebugMode {
				log.Printf("Error parsing event: %v", err)
			}
			continue
		}
		
		eventType, ok := event["type"].(string)
		if !ok {
			continue
		}
		
		if DebugMode {
			log.Printf("Event: %s", event)
		}
		
		switch eventType {
		case "challenge":
			// Challenge data is nested under "challenge" key
			challengeData, ok := event["challenge"].(map[string]interface{})
			if !ok {
				log.Printf("Error parsing challenge: challenge data not found")
				continue
			}
			
			// Convert to JSON and unmarshal into Challenge struct
			challengeJSON, err := json.Marshal(challengeData)
			if err != nil {
				log.Printf("Error marshaling challenge data: %v", err)
				continue
			}
			
			var challenge Challenge
			if err := json.Unmarshal(challengeJSON, &challenge); err != nil {
				log.Printf("Error parsing challenge: %v", err)
				continue
			}
			
			log.Printf("Received challenge: ID=%s, Status=%s", challenge.ID, challenge.Status)
			bot.handleChallenge(challenge)
			
		case "challengeCanceled":
			// A player canceled their challenge to us
			challengeData, ok := event["challenge"].(map[string]interface{})
			if ok {
				if challengeID, ok := challengeData["id"].(string); ok {
					log.Printf("Challenge %s was canceled", challengeID)
				}
			}
			
		case "challengeDeclined":
			// The opponent declined our challenge
			challengeData, ok := event["challenge"].(map[string]interface{})
			if ok {
				if challengeID, ok := challengeData["id"].(string); ok {
					log.Printf("Challenge %s was declined by opponent", challengeID)
				}
			}
			
		case "gameStart":
			// Start of a game - when stream opens, all current games are sent
			gameData, ok := event["game"].(map[string]interface{})
			if !ok {
				log.Printf("Error parsing gameStart: game data not found")
				continue
			}
			
			gameID, ok := gameData["gameId"].(string)
			if !ok {
				// Fallback to "id" field if "gameId" is not present
				gameID, ok = gameData["id"].(string)
				if !ok {
					log.Printf("Error parsing gameStart: missing game ID")
					continue
				}
			}
			
			log.Printf("Game started: %s", gameID)
			go bot.handleGame(gameID)
			
		case "gameFinish":
			// Completion of a game
			gameData, ok := event["game"].(map[string]interface{})
			if ok {
				gameID, ok := gameData["gameId"].(string)
				if !ok {
					gameID, _ = gameData["id"].(string)
				}
				if gameID != "" {
					log.Printf("Game finished: %s", gameID)
					// The game handler will detect the end from the game stream
				}
			}
			
		default:
			if DebugMode {
				log.Printf("Unknown event type: %s", eventType)
			}
		}
	}
	
	// If we get here, the scanner stopped (connection closed)
	scanErr := scanner.Err()
	if scanErr != nil {
		log.Printf("Scanner error: %v", scanErr)
		return scanErr
	}
	
	log.Println("Event stream closed by server")
	return fmt.Errorf("event stream closed")
}

func (bot *LichessBot) handleChallenge(challenge Challenge) {
	// Accept all challenges for now
	// You can add filtering logic here (e.g., only accept rated games, certain time controls, etc.)
	log.Printf("Handling challenge: ID=%s, Status=%s, Variant=%s", challenge.ID, challenge.Status, challenge.Variant.Key)
	
	if challenge.Status == "created" {
		log.Printf("Accepting challenge %s", challenge.ID)
		if err := bot.acceptChallenge(challenge.ID); err != nil {
			log.Printf("Error accepting challenge %s: %v", challenge.ID, err)
		} else {
			log.Printf("Successfully accepted challenge %s", challenge.ID)
		}
	} else {
		log.Printf("Challenge %s has status '%s', not accepting", challenge.ID, challenge.Status)
	}
}

func (bot *LichessBot) handleGame(gameID string) {
	url := fmt.Sprintf("%s/api/bot/game/stream/%s", LichessAPIBase, gameID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for game %s: %v", gameID, err)
		return
	}
	
	req.Header.Set("Authorization", "Bearer "+bot.token)
	
	// Use a client without timeout for long-running streams
	streamClient := &http.Client{
		Timeout: 0, // No timeout for streaming
	}
	resp, err := streamClient.Do(req)
	if err != nil {
		log.Printf("Error streaming game %s: %v", gameID, err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to stream game %s: %s - %s", gameID, resp.Status, string(body))
		return
	}
	
	var currentBoard BitBoard
	var isWhite bool
	var initialFenForGame string // Store initial FEN for this game
	
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		
		var gameData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &gameData); err != nil {
			log.Printf("Error parsing game data: %v", err)
			continue
		}
		
		eventType, ok := gameData["type"].(string)
		if !ok {
			continue
		}
		
		switch eventType {
		case "gameFull":
			var gameFull GameFull
			if err := json.Unmarshal([]byte(line), &gameFull); err != nil {
				log.Printf("Error parsing gameFull: %v", err)
				continue
			}
			
			// Determine if we're playing white or black
			// Cache user ID to avoid repeated API calls
			if bot.userID == "" {
				bot.userID = bot.getUserID()
				if bot.userID == "" {
					log.Printf("Could not determine bot user ID")
					return
				}
			}
			
			isWhite = gameFull.White.ID == bot.userID
			log.Printf("Game %s started. Playing as %s", gameID, map[bool]string{true: "White", false: "Black"}[isWhite])
			
			// Store initial FEN for this game (convert "startpos" to actual FEN)
			initialFenForGame = gameFull.InitialFen
			if initialFenForGame == "startpos" || initialFenForGame == "" {
				initialFenForGame = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
			}
			
			// Check if we can get the current FEN from gameFull.State.Fen
			// If not, convert initialFen (handling "startpos" special value)
			var currentFen string
			hasMoves := strings.TrimSpace(gameFull.State.Moves) != "" && strings.TrimSpace(gameFull.State.Moves) != "."
			
			if gameFull.State.Fen != "" {
				currentFen = gameFull.State.Fen
				log.Printf("Got FEN from gameFull.State: %s", currentFen)
			} else {
				currentFen = initialFenForGame
				if hasMoves {
					log.Printf("Game has moves: %s. Will reconstruct board from moves...", gameFull.State.Moves)
				} else {
					log.Printf("No moves yet. Using initial FEN: %s", initialFenForGame)
				}
			}
			
			// Set up board from FEN
			currentBoard = BoardFromFEN(currentFen)
			
			// If there are moves, apply them to get to current position
			if hasMoves {
				movesList := strings.Fields(gameFull.State.Moves)
				for _, uciMove := range movesList {
					newBoard, ok := ApplyUCIMove(currentBoard, uciMove)
					if !ok {
						log.Printf("Warning: Could not apply move %s. Current FEN: %s", uciMove, currentBoard.ToFEN())
						break
					}
					currentBoard = newBoard
				}
				log.Printf("Reconstructed board from moves. Current FEN: %s", currentBoard.ToFEN())
			}
			
			// Check if it's our turn
			if (isWhite && currentBoard.Turn() == White) || (!isWhite && currentBoard.Turn() == Black) {
				log.Printf("It's our turn! Making move from gameFull...")
				bot.makeBotMove(gameID, currentBoard)
			} else {
				log.Printf("Not our turn yet. Waiting for gameState event...")
			}
			
			// Check if game is already finished by parsing status
			var statusName string
			if len(gameFull.State.Status) > 0 {
				// Try to parse as object first
				var statusObj struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}
				if err := json.Unmarshal(gameFull.State.Status, &statusObj); err == nil {
					statusName = statusObj.Name
				} else {
					// Try to parse as string
					if err := json.Unmarshal(gameFull.State.Status, &statusName); err != nil {
						statusName = ""
					}
				}
			}
			
			if statusName != "" && statusName != "started" {
				log.Printf("Game %s already finished: %s", gameID, statusName)
				return
			}
			
			// Always wait for gameState event - it will have the current FEN and we'll make a move then
			
		case "gameState":
			log.Printf("Received gameState event for game %s", gameID)
			var gameState GameState
			if err := json.Unmarshal([]byte(line), &gameState); err != nil {
				log.Printf("Error parsing gameState: %v", err)
				log.Printf("Raw gameState line: %s", line)
				continue
			}
			
			// Update board from FEN if provided, otherwise reconstruct from moves
			if gameState.Fen != "" {
				log.Printf("Received gameState. FEN: %s", gameState.Fen)
				currentBoard = BoardFromFEN(gameState.Fen)
				log.Printf("Parsed board. Board FEN: %s", currentBoard.ToFEN())
			} else if gameState.Moves != "" {
				// Reconstruct board from moves
				log.Printf("gameState has no FEN, reconstructing from moves: %s", gameState.Moves)
				
				// Start from initial position
				if initialFenForGame == "" {
					initialFenForGame = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
				}
				currentBoard = BoardFromFEN(initialFenForGame)
				
				// Apply all moves
				movesList := strings.Fields(gameState.Moves)
				for _, uciMove := range movesList {
					newBoard, ok := ApplyUCIMove(currentBoard, uciMove)
					if !ok {
						log.Printf("Warning: Could not apply move %s. Current FEN: %s", uciMove, currentBoard.ToFEN())
						break
					}
					currentBoard = newBoard
				}
				log.Printf("Reconstructed board from moves. Current FEN: %s", currentBoard.ToFEN())
			} else {
				log.Printf("Warning: gameState event has no FEN and no moves. Raw data: %s", line)
				continue
			}
			
			// Check if game is over by parsing status (can be string or object)
			var statusName string
			if len(gameState.Status) > 0 {
				// Try to parse as object first
				var statusObj struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}
				if err := json.Unmarshal(gameState.Status, &statusObj); err == nil {
					statusName = statusObj.Name
				} else {
					// Try to parse as string
					if err := json.Unmarshal(gameState.Status, &statusName); err != nil {
						statusName = ""
					}
				}
			}
			
			if statusName != "" && statusName != "started" {
				log.Printf("Game %s ended: %s", gameID, statusName)
				if gameState.Winner != "" {
					if (isWhite && gameState.Winner == "white") || (!isWhite && gameState.Winner == "black") {
						log.Printf("Game %s: Bot won!", gameID)
					} else {
						log.Printf("Game %s: Bot lost", gameID)
					}
				}
				return
			}
			
			// If it's our turn, make a move
			if (isWhite && currentBoard.Turn() == White) || (!isWhite && currentBoard.Turn() == Black) {
				log.Printf("It's our turn! Current FEN: %s", currentBoard.ToFEN())
				bot.makeBotMove(gameID, currentBoard)
			} else {
				log.Printf("Not our turn. Waiting for opponent's move...")
			}
			
		case "chatLine":
			// Ignore chat messages
		case "opponentGone":
			// Opponent disconnected, wait for them to return
		default:
			// Check for draw/takeback offers in other event types
			if drawOffer, ok := gameData["drawOffer"].(bool); ok && drawOffer {
				bot.declineDraw(gameID)
			}
			if takebackOffer, ok := gameData["takebackOffer"].(bool); ok && takebackOffer {
				bot.declineTakeback(gameID)
			}
		}
	}
	
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading game stream for %s: %v", gameID, err)
	}
}

func (bot *LichessBot) makeBotMove(gameID string, board BitBoard) {
	log.Printf("Calculating best move for game %s...", gameID)
	
	// Get the best move
	evaluatedBoard := BestMoveWithoutTimeLimit(&board)
	
	log.Printf("Best move evaluation: %f", evaluatedBoard.eval)
	
	// Find the move that transforms current board to the best move board
	move, found := FindMoveBetweenBoards(&board, &evaluatedBoard.board)
	if !found {
		log.Printf("Could not find move between boards in game %s", gameID)
		log.Printf("From board FEN: %s", board.ToFEN())
		log.Printf("To board FEN: %s", evaluatedBoard.board.ToFEN())
		
		// Fallback: generate all legal moves and pick the first one
		// This shouldn't happen, but it's a safety net
		legalMoves, numMoves := GenerateLegalStates(&board)
		if numMoves > 0 {
			log.Printf("Fallback: using first legal move out of %d moves", numMoves)
			// Find the move that leads to the first legal state
			firstLegalBoard := legalMoves[0]
			move, found = FindMoveBetweenBoards(&board, &firstLegalBoard)
			if !found {
				log.Printf("Even fallback failed! This is unexpected.")
				return
			}
		} else {
			log.Printf("No legal moves available!")
			return
		}
	}
	
	// Convert to UCI format
	uciMove := move.ToUCI()
	
	log.Printf("Making move: %s in game %s", uciMove, gameID)
	
	// Make the move
	if err := bot.makeMove(gameID, uciMove); err != nil {
		log.Printf("Error making move %s in game %s: %v", uciMove, gameID, err)
	} else {
		log.Printf("Successfully made move %s", uciMove)
	}
}

func (bot *LichessBot) getUserID() string {
	url := fmt.Sprintf("%s/api/account", LichessAPIBase)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	
	req.Header.Set("Authorization", "Bearer "+bot.token)
	
	resp, err := bot.httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	
	var account struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return ""
	}
	
	return account.ID
}

func StartLichessBot(token string) {
	if token == "" {
		log.Fatal("LICHESS_TOKEN environment variable is not set")
	}
	
	bot := NewLichessBot(token)
	log.Println("Starting Lichess bot...")
	
	// Stream events in a loop (reconnect on error)
	// Per Lichess API docs: wait a full minute (60 seconds) after 429 errors
	// Also wait longer when stream closes to avoid appearing like polling
	consecutive429Errors := 0
	max429Wait := 5 // Maximum number of consecutive 429 errors before giving up
	
	for {
		if err := bot.streamEvents(); err != nil {
			// Check if it's a rate limit error (429)
			if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "Too Many Requests") {
				consecutive429Errors++
				if consecutive429Errors > max429Wait {
					log.Fatalf("Received %d consecutive 429 errors. Please wait at least 10 minutes before trying again. The rate limit may be in effect.", consecutive429Errors)
				}
				
				waitTime := 60 * time.Second
				if consecutive429Errors > 1 {
					// If we're still rate limited after waiting, wait longer
					waitTime = time.Duration(60+30*consecutive429Errors) * time.Second
					log.Printf("Still rate limited after %d attempts. Waiting %v before retrying...", consecutive429Errors, waitTime)
				} else {
					log.Printf("Rate limited (429). Waiting 60 seconds before resuming as per Lichess API requirements...")
				}
				time.Sleep(waitTime)
			} else {
				// Reset 429 counter on non-rate-limit errors
				consecutive429Errors = 0
				log.Printf("Error streaming events: %v. Waiting 10 seconds before reconnecting...", err)
				// Wait a bit longer on errors to avoid rapid reconnection
				time.Sleep(10 * time.Second)
			}
		} else {
			// Connection closed normally - reset counter and wait longer to avoid appearing like polling
			consecutive429Errors = 0
			log.Println("Event stream closed. Waiting 15 seconds before reconnecting...")
			time.Sleep(15 * time.Second)
		}
	}
}

