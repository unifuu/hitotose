package game

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Types
type Status string
type Genre string
type Platform string

const (
	// Status
	PLAYING = Status("Playing")
	PLAYED  = Status("Played")
	TO_PLAY = Status("ToPlay")

	// Platforms
	PC              = Platform("PC")
	NINTENDO_SWITCH = Platform("Nintendo Switch")
	PLAYSTATION     = Platform("PlayStation")
	XBOX            = Platform("Xbox")
	MOBILE          = Platform("Mobile")
)

// Return status
func Statuses() []Status {
	return []Status{
		PLAYING, PLAYED, TO_PLAY,
	}
}

// Badges
type Badge struct {
	// Status
	Played  int `json:"played"`
	Playing int `json:"playing"`
	ToPlay  int `json:"to_play"`

	// Platform
	AllPlatform    int `json:"all_platform"`
	PC             int `json:"pc"`
	PlayStation    int `json:"playstation"`
	NintendoSwitch int `json:"nintendo_switch"`
	XBox           int `json:"xbox"`
	Mobile         int `json:"mobile"`
}

// Game model
type Game struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Title      string             `json:"title" bson:"title"`
	Genre      string             `json:"genre" bson:"genre"`
	Platform   string             `json:"platform" bson:"platform"`
	Developer  string             `json:"developer" bson:"developer"`
	Publisher  string             `json:"publisher" bson:"publisher"`
	Status     Status             `json:"status" bson:"status"`
	PlayedTime int                `json:"played_time" bson:"played_time"`
	Rating     int                `json:"rating" bson:"rating"`
	CreatedAt  time.Time          `json:"-" bson:"created_at"`
	UpdatedAt  time.Time          `json:"-" bson:"updated_at"`
}

type StopWatch struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int       `json:"duration"`
	GameID    string    `json:"game_id"`
	GameTitle string    `json:"game_title"`
}

func NewStopWatch(id, title string) *StopWatch {
	return &StopWatch{
		StartTime: time.Time{},
		EndTime:   time.Time{},
		Duration:  0,
		GameID:    id,
		GameTitle: title,
	}
}

func (sw *StopWatch) Start() error {
	if len(sw.GameID) == 0 {
		return fmt.Errorf("game id cannot be empty")
	}
	sw.StartTime = time.Now()
	return nil
}

// Stop ends the stopwatch and returns the duration
func (sw *StopWatch) Stop() int {
	sw.EndTime = time.Now()
	sw.Duration = int(sw.EndTime.Sub(sw.StartTime).Minutes())
	return sw.Duration
}
