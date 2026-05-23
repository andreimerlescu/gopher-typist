package app

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// App wires together the router, widgets, and Gio window.
type App struct {
	router  *Router
	widgets *Widgets
	theme   *material.Theme
}

func newApp() *App {
	th := material.NewTheme()
	return &App{
		router:  NewRouter(),
		widgets: &Widgets{},
		theme:   th,
	}
}

// Run is the entry point called from main.go.
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

			// poll for async leaderboard load
			a.router.PollLeaderboard()

			// tick timers
			dt := float32(e.Now.Sub(e.Now).Seconds()) // always 0 — use predicted
			dt = gtx.Now.Sub(gtx.Now).Seconds()        // still 0
			_ = dt
			a.tick(gtx)

			// handle keyboard input
			a.handleInput(gtx, w)

			// draw current screen
			a.draw(gtx)

			e.Frame(gtx.Ops)
		}
	}
}
