// internal/leaderboard.go
package app

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

// magic is the 4-byte header identifying a rusty-typist leaderboard file.
// Kept identical to the Rust version so existing .bin files load correctly.
var magic = [4]byte{'R', 'T', 'L', 'B'}

const maxEntries = 20

type LeaderboardEntry struct {
	Name      string
	WPM       float64
	Accuracy  float64
	CPS       float64
	Level     uint8
	Timestamp uint64
}

// TimestampDisplay formats a UTC unix timestamp as "YYYY-MM-DD HH:MM".
// Pure stdlib, no time.Time — matches the Rust implementation exactly.
func (e LeaderboardEntry) TimestampDisplay() string {
	s := e.Timestamp
	const (
		secsPerMin  = 60
		secsPerHour = 3600
		secsPerDay  = 86400
	)

	days := s / secsPerDay
	tod := s % secsPerDay
	hh := tod / secsPerHour
	mm := (tod % secsPerHour) / secsPerMin

	year := uint64(1970)
	for {
		diy := uint64(365)
		if isLeap(year) {
			diy = 366
		}
		if days < diy {
			break
		}
		days -= diy
		year++
	}

	months := [12]uint64{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if isLeap(year) {
		months[1] = 29
	}
	month := uint64(1)
	for _, m := range months {
		if days < m {
			break
		}
		days -= m
		month++
	}
	day := days + 1

	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d", year, month, day, hh, mm)
}

func isLeap(y uint64) bool {
	return (y%4 == 0 && y%100 != 0) || y%400 == 0
}

type Leaderboard struct {
	Entries []LeaderboardEntry
	path    string
}

func EmptyLeaderboard() Leaderboard {
	return Leaderboard{path: configPath()}
}

func LoadLeaderboard() Leaderboard {
	path := configPath()
	lb := Leaderboard{path: path}

	data, err := os.ReadFile(path)
	if err != nil {
		return lb
	}

	entries, err := deserialize(data)
	if err != nil {
		return lb
	}

	lb.Entries = entries
	lb.sort()
	return lb
}

func (lb *Leaderboard) Save() {
	if dir := filepath.Dir(lb.path); dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}
	data, err := serialize(lb.Entries)
	if err != nil {
		return
	}
	_ = os.WriteFile(lb.path, data, 0644)
}

func (lb *Leaderboard) IsNewBest(wpm float64) bool {
	if len(lb.Entries) == 0 {
		return true
	}
	return wpm > lb.Entries[0].WPM
}

func (lb *Leaderboard) RankOf(wpm float64) int {
	rank := 1
	for _, e := range lb.Entries {
		if e.WPM >= wpm {
			rank++
		}
	}
	return rank
}

func (lb *Leaderboard) Insert(entry LeaderboardEntry) {
	lb.Entries = append(lb.Entries, entry)
	lb.sort()
	if len(lb.Entries) > maxEntries {
		lb.Entries = lb.Entries[:maxEntries]
	}
	lb.Save()
}

func (lb *Leaderboard) sort() {
	for i := 1; i < len(lb.Entries); i++ {
		for j := i; j > 0 && lb.Entries[j].WPM > lb.Entries[j-1].WPM; j-- {
			lb.Entries[j], lb.Entries[j-1] = lb.Entries[j-1], lb.Entries[j]
		}
	}
}

// configPath returns ~/.config/rusty-typist/leaderboard.bin with fallback.
// Uses the same path as the Rust version so scores carry over.
func configPath() string {
	home, err := os.UserHomeDir()
	if err == nil {
		dir := filepath.Join(home, ".config", "rusty-typist")
		if os.MkdirAll(dir, 0755) == nil {
			probe := filepath.Join(dir, ".write_probe")
			if f, err := os.Create(probe); err == nil {
				f.Close()
				os.Remove(probe)
				return filepath.Join(dir, "leaderboard.bin")
			}
		}
	}

	// fallback: next to the binary
	exe, err := os.Executable()
	if err == nil {
		return filepath.Join(filepath.Dir(exe), "rusty-typist-leaderboard.bin")
	}
	return "rusty-typist-leaderboard.bin"
}

// serialize writes entries in the RTLB binary format.
// Format: [RTLB magic][u32 count][entries...]
// Each entry: [u8 name_len][name bytes][f64 wpm][f64 accuracy][f64 cps][u8 level][u64 timestamp]
func serialize(entries []LeaderboardEntry) ([]byte, error) {
	buf := make([]byte, 0, 256)
	buf = append(buf, magic[:]...)

	count := make([]byte, 4)
	binary.LittleEndian.PutUint32(count, uint32(len(entries)))
	buf = append(buf, count...)

	tmp := make([]byte, 8)
	for _, e := range entries {
		name := []byte(e.Name)
		if len(name) > 255 {
			name = name[:255]
		}
		buf = append(buf, byte(len(name)))
		buf = append(buf, name...)

		binary.LittleEndian.PutUint64(tmp, floatBits(e.WPM))
		buf = append(buf, tmp...)
		binary.LittleEndian.PutUint64(tmp, floatBits(e.Accuracy))
		buf = append(buf, tmp...)
		binary.LittleEndian.PutUint64(tmp, floatBits(e.CPS))
		buf = append(buf, tmp...)

		buf = append(buf, e.Level)

		binary.LittleEndian.PutUint64(tmp, e.Timestamp)
		buf = append(buf, tmp...)
	}
	return buf, nil
}

func deserialize(data []byte) ([]LeaderboardEntry, error) {
	if len(data) < 8 {
		return nil, io.ErrUnexpectedEOF
	}
	if data[0] != 'R' || data[1] != 'T' || data[2] != 'L' || data[3] != 'B' {
		return nil, fmt.Errorf("bad magic")
	}

	count := int(binary.LittleEndian.Uint32(data[4:8]))
	entries := make([]LeaderboardEntry, 0, count)
	pos := 8

	for range count {
		if pos >= len(data) {
			break
		}
		nameLen := int(data[pos])
		pos++
		if pos+nameLen > len(data) {
			break
		}
		name := string(data[pos : pos+nameLen])
		pos += nameLen

		if pos+33 > len(data) {
			break
		}
		wpm := bitsFloat(binary.LittleEndian.Uint64(data[pos:]))
		pos += 8
		accuracy := bitsFloat(binary.LittleEndian.Uint64(data[pos:]))
		pos += 8
		cps := bitsFloat(binary.LittleEndian.Uint64(data[pos:]))
		pos += 8
		level := data[pos]
		pos++
		timestamp := binary.LittleEndian.Uint64(data[pos:])
		pos += 8

		entries = append(entries, LeaderboardEntry{
			Name: name, WPM: wpm, Accuracy: accuracy,
			CPS: cps, Level: level, Timestamp: timestamp,
		})
	}
	return entries, nil
}

func floatBits(f float64) uint64 { 
	return math.Float64bits(f) 
}
func bitsFloat(u uint64) float64 { 
	return math.Float64frombits(u) 
}

