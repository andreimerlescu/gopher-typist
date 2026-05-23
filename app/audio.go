// internal/audio.go
package app

import (
	"os"
	"path/filepath"
)

type audioMsg struct {
	kind string // "sfx", "bg_start", "bg_stop"
	path string
}

type AudioManager struct {
	ch chan audioMsg
}

// NewAudioManager spawns the audio goroutine and returns immediately.
// Returns nil if audio initialization fails — callers must nil-check.
func NewAudioManager() *AudioManager {
	ch := make(chan audioMsg, 16)
	am := &AudioManager{ch: ch}
	go am.run()
	return am
}

func (a *AudioManager) PlaySFX(path string) {
	if a == nil {
		return
	}
	a.ch <- audioMsg{kind: "sfx", path: path}
}

func (a *AudioManager) StartBackground(path string) {
	if a == nil {
		return
	}
	a.ch <- audioMsg{kind: "bg_start", path: path}
}

func (a *AudioManager) StopBackground() {
	if a == nil {
		return
	}
	a.ch <- audioMsg{kind: "bg_stop"}
}

// Asset resolves a WAV filename to an absolute path.
// Checks relative to CWD first, then relative to the executable.
func Asset(name string) string {
	candidates := []string{
		filepath.Join("assets", "wav", name),
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates,
			filepath.Join(filepath.Dir(exe), "assets", "wav", name),
		)
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	// return the first candidate as default even if it doesn't exist
	return candidates[0]
}
