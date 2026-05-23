package app

import (
	"testing"
	"time"
)

func TestNewGameInitializesCorrectly(t *testing.T) {
	g := NewGame(ModeQuick, 0)

	if g.Cursor != 0 {
		t.Errorf("want cursor=0, got %d", g.Cursor)
	}
	if len(g.Text) == 0 {
		t.Error("text should not be empty")
	}
	if len(g.Correctness) != len(g.Text) {
		t.Errorf("correctness len %d != text len %d", len(g.Correctness), len(g.Text))
	}
	if g.StartTime != nil {
		t.Error("start time should be nil before first keypress")
	}
}

func TestHandleCharCorrectAdvancesCursor(t *testing.T) {
	g := NewGame(ModeQuick, 0)
	first := g.Text[0]

	event := g.HandleChar(first)

	if event != GameEventCorrect && event != GameEventLevelComplete {
		t.Errorf("want Correct or LevelComplete, got %d", event)
	}
	if g.Cursor != 1 {
		t.Errorf("want cursor=1, got %d", g.Cursor)
	}
	if g.Correctness[0] != CorrectnessCorrect {
		t.Error("want correctness[0]=Correct")
	}
}

func TestHandleCharWrongQuickModeDoesNotAdvance(t *testing.T) {
	g := NewGame(ModeQuick, 0)

	// find a char that is definitely not the first char
	wrong := rune('!')
	if g.Text[0] == wrong {
		wrong = rune('@')
	}

	event := g.HandleChar(wrong)

	if event != GameEventWrong {
		t.Errorf("want Wrong, got %d", event)
	}
	if g.Cursor != 0 {
		t.Errorf("quick mode: wrong key should not advance cursor, got %d", g.Cursor)
	}
	if g.ShakeTimer <= 0 {
		t.Error("shake timer should be set on wrong key")
	}
}

func TestHandleCharWrongHardModeAdvances(t *testing.T) {
	g := NewGame(ModeHard, 0)

	wrong := rune('!')
	if g.Text[0] == wrong {
		wrong = rune('@')
	}

	event := g.HandleChar(wrong)

	if event != GameEventWrong {
		t.Errorf("want Wrong, got %d", event)
	}
	if g.Cursor != 1 {
		t.Errorf("hard mode: wrong key should advance cursor, got %d", g.Cursor)
	}
	if g.Correctness[0] != CorrectnessWrong {
		t.Error("want correctness[0]=Wrong")
	}
}

func TestStartTimeSetOnFirstKeypress(t *testing.T) {
	g := NewGame(ModeQuick, 0)
	if g.StartTime != nil {
		t.Fatal("start time should be nil before first key")
	}
	g.HandleChar(g.Text[0])
	if g.StartTime == nil {
		t.Error("start time should be set after first key")
	}
}

func TestTickTimerQuickModeNeverExpires(t *testing.T) {
	g := NewGame(ModeQuick, 0)
	now := time.Now()
	g.StartTime = &now

	expired := g.TickTimer(9999.0)
	if expired {
		t.Error("quick mode timer should never expire")
	}
}

func TestTickTimerHardModeExpires(t *testing.T) {
	g := NewGame(ModeHard, 0)
	now := time.Now()
	g.StartTime = &now
	g.RemainingeSecs = 1.0

	expired := g.TickTimer(2.0)
	if !expired {
		t.Error("hard mode timer should expire when dt > remaining")
	}
	if g.RemainingeSecs != 0 {
		t.Errorf("remaining should clamp to 0, got %f", g.RemainingeSecs)
	}
}

func TestTickTimerRequiresStartTime(t *testing.T) {
	g := NewGame(ModeHard, 0)
	// StartTime is nil — timer should not tick
	expired := g.TickTimer(999.0)
	if expired {
		t.Error("timer should not expire before first keypress")
	}
}

func TestScoreCalculation(t *testing.T) {
	g := NewGame(ModeQuick, 0)

	// type the full level correctly
	for i, ch := range g.Text {
		event := g.HandleChar(ch)
		if event == GameEventLevelComplete && i < len(g.Text)-1 {
			t.Fatal("level completed too early")
		}
	}

	score := g.Score()

	if score.WPM <= 0 {
		t.Errorf("want WPM > 0, got %f", score.WPM)
	}
	if score.Accuracy != 100.0 {
		t.Errorf("want 100%% accuracy for all-correct run, got %f", score.Accuracy)
	}
	if score.CPS <= 0 {
		t.Errorf("want CPS > 0, got %f", score.CPS)
	}
	if score.Level != 0 {
		t.Errorf("want level=0, got %d", score.Level)
	}
}

func TestScoreAccuracyWithMistakes(t *testing.T) {
	g := NewGame(ModeHard, 0) // hard mode advances on wrong too

	// type one wrong then one correct
	wrong := rune('!')
	if g.Text[0] == wrong {
		wrong = rune('@')
	}
	g.HandleChar(wrong)  // wrong
	g.HandleChar(g.Text[1]) // correct (cursor is now at 1 after hard-mode advance)

	score := g.Score()
	if score.Accuracy >= 100.0 {
		t.Error("accuracy should be less than 100 after a mistake")
	}
}

func TestRandomizePositionStaysInBounds(t *testing.T) {
	g := NewGame(ModeHard, 0)

	for range 1000 {
		g.randomizePosition()
		if g.Position.X < 0.15 || g.Position.X > 0.85 {
			t.Errorf("X out of bounds: %f", g.Position.X)
		}
		if g.Position.Y < 0.20 || g.Position.Y > 0.80 {
			t.Errorf("Y out of bounds: %f", g.Position.Y)
		}
	}
}

func TestHardModeRemainingSecsCalculation(t *testing.T) {
	g := NewGame(ModeHard, 0)
	wordCount := float64(len(g.Text)) / 5.0
	expected := float32(wordCount/HardTargetWPM*60.0) + hardGraceSeconds

	if g.RemainingeSecs != expected {
		t.Errorf("want %.2f remaining secs, got %.2f", expected, g.RemainingeSecs)
	}
}

func TestLevelCompleteFiredWhenCursorReachesEnd(t *testing.T) {
	g := NewGame(ModeQuick, 0)
	var lastEvent GameEvent
	for _, ch := range g.Text {
		lastEvent = g.HandleChar(ch)
	}
	if lastEvent != GameEventLevelComplete {
		t.Errorf("want LevelComplete on last char, got %d", lastEvent)
	}
}
