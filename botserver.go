package main

// EventType enum
type EventType string

// Event struct
type Event struct {
	typee     string    `json:"type"` // "gameStart", "gameFinish", "challenge", "challengeCanceled", "challengeDeclined"
	game      Game      `json:"game"`
	challenge Challenge `json:"challenge"`
}

// Game struct
type Game struct {
	GameId      string   `json:"gameId"`
	FullId      string   `json:"fullId"`
	Color       string   `json:"color"` // "white", "black"
	Fen         string   `json:"fen"`
	HasMoved    bool     `json:"hasMoved"`
	IsMyTurn    bool     `json:"isMyTurn"`
	LastMove    string   `json:"lastMove"`
	Opponent    Opponent `json:"opponent"`
	Perf        string   `json:"perf"`
	Rated       bool     `json:"rated"`
	SecondsLeft int      `json:"secondsLeft"`
	Source      string   `json:"source"`
	Status      Status   `json:"status"`
	Speed       string   `json:"speed"`
	Variant     Variant  `json:"variant"`
	Compat      Compat   `json:"compat"`
	Winner      string   `json:"winner"`
	RatingDiff  int      `json:"ratingDiff"`
	Id          string   `json:"id"`
}

type Challenge struct {
	Id          string      `json:"id"`
	Url         string      `json:"url"`
	Status      string      `json:"status"` // "created", "canceled", "declined"
	Challenger  Challenger  `json:"challenger"`
	DestUser    DestUser    `json:"destUser"`
	Variant     Variant     `json:"variant"`
	Rated       bool        `json:"rated"`
	Speed       string      `json:"speed"`
	TimeControl TimeControl `json:"timeControl"`
	Color       string      `json:"color"`      // "white", "black", "random"
	FinalColor  string      `json:"finalColor"` // "white", "black"
	Perf        Perf        `json:"perf"`
	Direction   string      `json:"direction"`
}

type Challenger struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Rating      int    `json:"rating"`
	Title       string `json:"title"`
	Provisional bool   `json:"provisional"`
	Online      bool   `json:"online"`
	Lag         int    `json:"lag"`
}

type DestUser struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Rating      int    `json:"rating"`
	Title       string `json:"title"`
	Provisional bool   `json:"provisional"`
	Online      bool   `json:"online"`
	Lag         int    `json:"lag"`
}

type TimeControl struct {
	Type      string `json:"type"`
	Limit     int    `json:"limit"`
	Increment int    `json:"increment"`
	Show      string `json:"show"`
}

type Perf struct {
	Icon string `json:"icon"`
	Name string `json:"name"`
}

// Opponent struct
type Opponent struct {
	Id       string `json:"id"`
	Rating   int    `json:"rating"`
	Username string `json:"username"`
}

// Status struct
type Status struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// Variant struct
type Variant struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// Compat struct
type Compat struct {
	Bot   bool `json:"bot"`
	Board bool `json:"board"`
}
