package app

import "gioui.org/widget"

// Widgets holds all stateful Gio widgets for the lifetime of the app.
// Declared once in App, never re-created.
type Widgets struct {
	// menu
	QuickBtn widget.Clickable
	HardBtn  widget.Clickable

	// name entry
	NameEditor widget.Editor
	SubmitBtn  widget.Clickable

	// score screen
	NextLevelBtn widget.Clickable
	PlayAgainBtn widget.Clickable
	MenuBtn      widget.Clickable

	// timeout screen
	RetryBtn    widget.Clickable
	TimeoutMenu widget.Clickable
}
