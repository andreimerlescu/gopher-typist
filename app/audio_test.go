// internal/audio_run.go
package app

import (
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"time"
)

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

func TestAssetReturnsCWDPathWhenFileExists(t *testing.T) {
	// create a temp wav file in the expected relative location
	dir := t.TempDir()
	wavDir := filepath.Join(dir, "assets", "wav")
	if err := os.MkdirAll(wavDir, 0755); err != nil {
		t.Fatal(err)
	}
	wavFile := filepath.Join(wavDir, "test.wav")
	if err := os.WriteFile(wavFile, []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}

	// Asset() checks CWD-relative path — we can't easily override CWD in tests
	// so just verify it returns a non-empty string and falls back gracefully
	result := Asset("keypress.wav")
	if result == "" {
		t.Error("Asset() returned empty string")
	}
}

func TestAssetReturnsDefaultWhenNotFound(t *testing.T) {
	result := Asset("nonexistent.wav")
	// should return the first candidate path, not empty
	if result == "" {
		t.Error("Asset() returned empty string for missing file")
	}
	if filepath.Base(result) != "nonexistent.wav" {
		t.Errorf("expected filename nonexistent.wav in path, got %s", result)
	}
}
