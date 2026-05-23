package app

import (
	"math/rand"
	"time"
)

// ── Constants ─────────────────────────────────────────────────────────────────

const (
	HardTargetWPM = 20.0
	CharFontSize  = 144.0
	WindowRadius  = 3

	hardGraceSeconds = 15.0
)

var Levels = []string{
	"red orange yellow green blue indigo violet are the colors of the rainbow.",
	"yeshua loves me. yeshua has plans to prosper me. i love yeshua.",
	"i am beautiful. i am worthy. i am brave. i am confident. i am kind.",
	"solutions present themselves daily to me. money flows to me when I need it.",
	"no weapon formed against me prospers. yeshua protects me.",
	"no conspiracy waged against me prospers. yeshua prospers me.",
	"i the lord you god will deliver you from the evil one. yeshua saves me.",
}

// ── Types ─────────────────────────────────────────────────────────────────────

type Mode int

const (
	ModeQuick Mode = iota
	ModeHard
)

type Correctness int

const (
	CorrectnessUntried Correctness = iota
	CorrectnessCorrect
	CorrectnessWrong
)

type LevelScore struct {
	Level    int
	Accuracy float64
	WPM      float64
	CPS      float64
}

type Position struct {
	X, Y float32
}

// ── Game ──────────────────────────────────────────────────────────────────────

type Game struct {
	Mode  Mode
	Level int

	Text        []rune
	Cursor      int
	Correctness []Correctness
	TotalTries  int

	StartTime     *time.Time
	RemainingeSecs float32
	ShakeTimer    float32
	Position      Position

	rng *rand.Rand
}

func NewGame(mode Mode, level int) *Game {
	g := &Game{
		Mode:  mode,
		Level: level,
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	g.reset()
	return g
}

func (g *Game) reset() {
	text := Levels[g.Level]
	g.Text = []rune(text)
	g.Cursor = 0
	g.Correctness = make([]Correctness, len(g.Text))
	g.TotalTries = 0
	g.StartTime = nil
	g.ShakeTimer = 0

	wordCount := float64(len(g.Text)) / 5.0
	if g.Mode == ModeHard {
		g.RemainingeSecs = float32(wordCount/HardTargetWPM*60.0) + hardGraceSeconds
	} else {
		g.RemainingeSecs = 0 // unused in Quick mode
	}

	g.randomizePosition()
}

func (g *Game) randomizePosition() {
	g.Position = Position{
		X: 0.15 + g.rng.Float32()*0.70,
		Y: 0.20 + g.rng.Float32()*0.60,
	}
}

// HandleChar processes a single keystroke. Returns a GameEvent indicating
// what happened so the caller can trigger audio/navigation without game.go
// importing either.
func (g *Game) HandleChar(ch rune) GameEvent {
	if g.StartTime == nil {
		now := time.Now()
		g.StartTime = &now
	}

	expected := g.Text[g.Cursor]
	g.TotalTries++

	if ch == expected {
		g.Correctness[g.Cursor] = CorrectnessCorrect
		g.Cursor++
		if g.Mode == ModeHard {
			g.randomizePosition()
		}
		if g.Cursor == len(g.Text) {
			return GameEventLevelComplete
		}
		return GameEventCorrect
	}

	// wrong key
	g.ShakeTimer = 0.5
	if g.Mode == ModeHard {
		g.Correctness[g.Cursor] = CorrectnessWrong
		g.Cursor++
		g.randomizePosition()
		if g.Cursor == len(g.Text) {
			return GameEventLevelComplete
		}
	}
	return GameEventWrong
}

// TickTimer advances the Hard mode countdown by dt seconds.
// Returns true if time has expired.
func (g *Game) TickTimer(dt float32) bool {
	if g.Mode != ModeHard || g.StartTime == nil {
		return false
	}
	g.RemainingeSecs -= dt
	if g.RemainingeSecs <= 0 {
		g.RemainingeSecs = 0
		return true
	}
	return false
}

// TickShake advances the shake animation timer.
func (g *Game) TickShake(dt float32) {
	if g.ShakeTimer > 0 {
		g.ShakeTimer -= dt
		if g.ShakeTimer < 0 {
			g.ShakeTimer = 0
		}
	}
}

// Score calculates the final score for the completed level.
func (g *Game) Score() LevelScore {
	elapsed := 1.0 // default to avoid div/0
	if g.StartTime != nil {
		elapsed = time.Since(*g.StartTime).Seconds()
		if elapsed < 0.001 {
			elapsed = 0.001
		}
	}

	correct := 0
	for _, c := range g.Correctness {
		if c == CorrectnessCorrect {
			correct++
		}
	}

	return LevelScore{
		Level:    g.Level,
		WPM:      (float64(len(g.Text)) / 5.0) / (elapsed / 60.0),
		Accuracy: float64(correct) / float64(max(g.TotalTries, 1)) * 100.0,
		CPS:      float64(len(g.Text)) / elapsed,
	}
}

// HardTimerTotal returns the total time allotted for the current level in Hard mode.
// Used by the draw layer to render a progress bar if desired.
func (g *Game) HardTimerTotal() float32 {
	wordCount := float64(len(g.Text)) / 5.0
	return float32(wordCount/HardTargetWPM*60.0) + hardGraceSeconds
}

// ── GameEvent ─────────────────────────────────────────────────────────────────

type GameEvent int

const (
	GameEventNone GameEvent = iota
	GameEventCorrect
	GameEventWrong
	GameEventLevelComplete
)
