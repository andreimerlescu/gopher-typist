// internal/leaderboard_test.go
package app

import (
	"os"
	"testing"
)

func TestSerializeDeserializeRoundTrip(t *testing.T) {
	entries := []LeaderboardEntry{
		{Name: "Alice", WPM: 80.5, Accuracy: 98.2, CPS: 6.7, Level: 2, Timestamp: 1700000000},
		{Name: "Bob", WPM: 65.0, Accuracy: 91.0, CPS: 5.4, Level: 1, Timestamp: 1700000001},
	}

	data, err := serialize(entries)
	if err != nil {
		t.Fatalf("serialize failed: %v", err)
	}

	got, err := deserialize(data)
	if err != nil {
		t.Fatalf("deserialize failed: %v", err)
	}

	if len(got) != len(entries) {
		t.Fatalf("want %d entries, got %d", len(entries), len(got))
	}
	for i, e := range entries {
		g := got[i]
		if g.Name != e.Name || g.WPM != e.WPM || g.Accuracy != e.Accuracy ||
			g.CPS != e.CPS || g.Level != e.Level || g.Timestamp != e.Timestamp {
			t.Errorf("entry %d mismatch: want %+v got %+v", i, e, g)
		}
	}
}

func TestBadMagicReturnsError(t *testing.T) {
	_, err := deserialize([]byte("XXXX\x01\x00\x00\x00"))
	if err == nil {
		t.Fatal("expected error for bad magic, got nil")
	}
}

func TestIsNewBest(t *testing.T) {
	lb := EmptyLeaderboard()
	if !lb.IsNewBest(1.0) {
		t.Fatal("empty leaderboard should always be new best")
	}

	lb.Entries = []LeaderboardEntry{{WPM: 50.0}}
	if !lb.IsNewBest(51.0) {
		t.Fatal("51 should beat 50")
	}
	if lb.IsNewBest(49.0) {
		t.Fatal("49 should not beat 50")
	}
}

func TestRankOf(t *testing.T) {
	lb := EmptyLeaderboard()
	lb.Entries = []LeaderboardEntry{
		{WPM: 100}, {WPM: 80}, {WPM: 60},
	}
	if r := lb.RankOf(90.0); r != 2 {
		t.Errorf("want rank 2, got %d", r)
	}
	if r := lb.RankOf(50.0); r != 4 {
		t.Errorf("want rank 4, got %d", r)
	}
}

func TestInsertTruncatesAt20(t *testing.T) {
	lb := EmptyLeaderboard()
	lb.path = filepath.Join(t.TempDir(), "test.bin")

	for i := range 25 {
		lb.Insert(LeaderboardEntry{Name: "x", WPM: float64(i)})
	}
	if len(lb.Entries) > maxEntries {
		t.Errorf("want max %d entries, got %d", maxEntries, len(lb.Entries))
	}
}

func TestInsertSortsByWPMDescending(t *testing.T) {
	lb := EmptyLeaderboard()
	lb.path = filepath.Join(t.TempDir(), "test.bin")

	lb.Insert(LeaderboardEntry{Name: "slow", WPM: 30.0})
	lb.Insert(LeaderboardEntry{Name: "fast", WPM: 90.0})
	lb.Insert(LeaderboardEntry{Name: "mid", WPM: 60.0})

	if lb.Entries[0].WPM != 90.0 {
		t.Errorf("want first entry WPM=90, got %.1f", lb.Entries[0].WPM)
	}
}

func TestTimestampDisplay(t *testing.T) {
	e := LeaderboardEntry{Timestamp: 0}
	if got := e.TimestampDisplay(); got != "1970-01-01 00:00" {
		t.Errorf("want 1970-01-01 00:00, got %s", got)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "leaderboard.bin")

	lb := Leaderboard{path: path}
	lb.Insert(LeaderboardEntry{Name: "Test", WPM: 77.7, Accuracy: 95.0, CPS: 6.0, Level: 1, Timestamp: 1234567890})

	// read it back using deserialize directly to avoid configPath() in LoadLeaderboard
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not written: %v", err)
	}
	entries, err := deserialize(data)
	if err != nil {
		t.Fatalf("deserialize failed: %v", err)
	}
	if len(entries) != 1 || entries[0].Name != "Test" {
		t.Errorf("unexpected entries: %+v", entries)
	}
}
