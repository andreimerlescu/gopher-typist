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
