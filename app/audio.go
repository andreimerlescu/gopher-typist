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

func (a *AudioManager) run() {
	var bgCtrl *beep.Ctrl
	initialized := false
	var format beep.Format

	initSpeaker := func(f beep.Format) {
		if !initialized {
			speaker.Init(f.SampleRate, f.SampleRate.N(time.Second/10))
			format = f
			initialized = true
		}
	}

	for msg := range a.ch {
		switch msg.kind {
		case "sfx":
			f, err := os.Open(msg.path)
			if err != nil {
				continue
			}
			streamer, fmt, err := wav.Decode(f)
			if err != nil {
				f.Close()
				continue
			}
			initSpeaker(fmt)
			speaker.Play(beep.Seq(streamer, beep.Callback(func() {
				f.Close()
			})))

		case "bg_start":
			f, err := os.Open(msg.path)
			if err != nil {
				continue
			}
			streamer, fmt, err := wav.Decode(f)
			if err != nil {
				f.Close()
				continue
			}
			initSpeaker(fmt)
			loop := beep.Loop(-1, streamer)
			ctrl := &beep.Ctrl{Streamer: loop, Paused: false}
			if bgCtrl != nil {
				speaker.Lock()
				bgCtrl.Streamer = nil
				speaker.Unlock()
			}
			bgCtrl = ctrl
			vol := &beep.Volume{
				Streamer: ctrl,
				Base:     2,
				Volume:   -1.5, // roughly 0.35 volume like the Rust version
				Silent:   false,
			}
			_ = format // suppress unused warning if sfx ran first
			speaker.Play(vol)

		case "bg_stop":
			if bgCtrl != nil {
				speaker.Lock()
				bgCtrl.Streamer = nil
				speaker.Unlock()
				bgCtrl = nil
			}
		}
	}
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
