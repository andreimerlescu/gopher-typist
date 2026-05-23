package app

import "time"

// ── Screen ────────────────────────────────────────────────────────────────────

type Screen int

const (
	ScreenMenu Screen = iota
	ScreenPlaying
	ScreenNameEntry
	ScreenScore
	ScreenTimeout
)

// ── NameEntryData and ScoreData carry per-screen state ────────────────────────

type NameEntryData struct {
	Score      LevelScore
	IsNewBest  bool
}

type ScoreData struct {
	Score      LevelScore
	PlayerRank int
}

// ── Router ────────────────────────────────────────────────────────────────────

// Router owns all navigation state and the current game instance.
// Draw functions read from it; App methods mutate it.
type Router struct {
	Screen      Screen
	NameEntry   NameEntryData
	Score       ScoreData

	Game        *Game
	Mode        Mode
	CurrentLevel int

	Leaderboard    Leaderboard
	leaderboardCh  chan Leaderboard

	NameInput   string
	Audio       *AudioManager
}

func NewRouter() *Router {
	// load leaderboard off the main goroutine so startup is instant
	ch := make(chan Leaderboard, 1)
	go func() { ch <- LoadLeaderboard() }()

	return &Router{
		Screen:        ScreenMenu,
		leaderboardCh: ch,
		Leaderboard:   EmptyLeaderboard(),
		Audio:         NewAudioManager(),
	}
}

// PollLeaderboard receives the leaderboard from the background load if ready.
// Call once per frame — returns immediately if not ready yet.
func (r *Router) PollLeaderboard() {
	if r.leaderboardCh == nil {
		return
	}
	select {
	case lb := <-r.leaderboardCh:
		r.Leaderboard = lb
		r.leaderboardCh = nil
	default:
	}
}

// ── Navigation ────────────────────────────────────────────────────────────────

func (r *Router) GoToMenu() {
	r.Audio.StopBackground()
	r.Screen = ScreenMenu
	r.Game = nil
}

func (r *Router) StartGame(mode Mode) {
	r.Mode = mode
	r.CurrentLevel = 0
	r.Audio.StartBackground(Asset("background.wav"))
	r.startLevel()
}

func (r *Router) startLevel() {
	r.Game = NewGame(r.Mode, r.CurrentLevel)
	r.Screen = ScreenPlaying
}

func (r *Router) NextLevel() {
	r.CurrentLevel++
	if r.CurrentLevel >= len(Levels) {
		r.GoToMenu()
		return
	}
	r.startLevel()
}

func (r *Router) FinishLevel() {
	score := r.Game.Score()
	isNewBest := r.Leaderboard.IsNewBest(score.WPM)

	if isNewBest {
		r.Audio.PlaySFX(Asset("new_best.wav"))
	} else {
		r.Audio.PlaySFX(Asset("level_complete.wav"))
	}

	r.CurrentLevel++
	r.NameInput = ""
	r.Screen = ScreenNameEntry
	r.NameEntry = NameEntryData{Score: score, IsNewBest: isNewBest}
}

func (r *Router) SubmitName() {
	name := r.NameInput
	if len([]rune(name)) == 0 {
		name = "Anonymous"
	}

	entry := LeaderboardEntry{
		Name:      name,
		WPM:       r.NameEntry.Score.WPM,
		Accuracy:  r.NameEntry.Score.Accuracy,
		CPS:       r.NameEntry.Score.CPS,
		Level:     uint8(r.NameEntry.Score.Level),
		Timestamp: uint64(time.Now().Unix()),
	}

	r.Leaderboard.Insert(entry)
	rank := r.Leaderboard.RankOf(r.NameEntry.Score.WPM)
	r.Screen = ScreenScore
	r.Score = ScoreData{Score: r.NameEntry.Score, PlayerRank: rank}
}

func (r *Router) Timeout() {
	r.Audio.PlaySFX(Asset("timeout.wav"))
	r.Screen = ScreenTimeout
}

func (r *Router) HasNextLevel() bool {
	return r.CurrentLevel < len(Levels)
}
