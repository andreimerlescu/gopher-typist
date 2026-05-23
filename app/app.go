package app

import (
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type App struct {
	router    *Router
	widgets   *Widgets
	theme     *material.Theme
	lastFrame time.Time
}

func newApp() *App {
	return &App{
		router:    NewRouter(),
		widgets:   &Widgets{},
		theme:     material.NewTheme(),
		lastFrame: time.Now(),
	}
}

func Run() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("Gopher Typist"),
			app.Size(unit.Dp(900), unit.Dp(600)),
		)
		a := newApp()
		if err := a.loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func (a *App) loop(w *app.Window) error {
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// dt in seconds since last frame
			now := time.Now()
			dt := float32(now.Sub(a.lastFrame).Seconds())
			a.lastFrame = now

			// non-blocking leaderboard poll
			a.router.PollLeaderboard()

			// tick game timers and shake
			if a.router.Screen == ScreenPlaying && a.router.Game != nil {
				a.router.Game.TickShake(dt)
				if a.router.Game.TickTimer(dt) {
					a.router.Timeout()
				}
			}

			// keyboard input
			a.handleInput(gtx, w)

			// draw
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return a.draw(gtx)
			})

			// request next frame so timers keep ticking
			w.Invalidate()

			e.Frame(gtx.Ops)
		}
	}
}

func (a *App) handleInput(gtx layout.Context, w *app.Window) {
	for {
		ev, ok := gtx.Event(key.Filter{})
		if !ok {
			break
		}
		ke, ok := ev.(key.Event)
		if !ok || ke.State != key.Press {
			continue
		}

		// Escape always goes to menu
		if ke.Name == key.NameEscape {
			a.router.GoToMenu()
			continue
		}

		// Enter submits name on name entry screen
		if ke.Name == key.NameReturn && a.router.Screen == ScreenNameEntry {
			a.router.NameInput = a.widgets.NameEditor.Text()
			a.router.SubmitName()
			a.widgets.NameEditor.SetText("")
			continue
		}

		// typing during game — handled via key.Filter for runes
		if a.router.Screen == ScreenPlaying && a.router.Game != nil {
			if len(ke.Name) == 1 {
				event := a.router.Game.HandleChar(rune(ke.Name[0]))
				switch event {
				case GameEventCorrect:
					a.router.Audio.PlaySFX(Asset("keypress.wav"))
				case GameEventWrong:
					a.router.Audio.PlaySFX(Asset("keypress_wrong.wav"))
				case GameEventLevelComplete:
					a.router.FinishLevel()
				}
			}
		}
	}
}

func (a *App) draw(gtx layout.Context) layout.Dimensions {
	switch a.router.Screen {
	case ScreenMenu:
		return DrawMenu(gtx, a.theme, a.router, a.widgets)
	case ScreenPlaying:
		return DrawPlaying(gtx, a.theme, a.router)
	case ScreenNameEntry:
		return DrawNameEntry(gtx, a.theme, a.router, a.widgets)
	case ScreenScore:
		return DrawScore(gtx, a.theme, a.router, a.widgets)
	case ScreenTimeout:
		return DrawTimeout(gtx, a.theme, a.router, a.widgets)
	default:
		return layout.Dimensions{}
	}
}
