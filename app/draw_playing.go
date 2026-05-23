package app

import (
	"fmt"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func DrawPlaying(gtx layout.Context, th *material.Theme, r *Router) layout.Dimensions {
	g := r.Game
	if g == nil {
		return layout.Dimensions{}
	}

	rect := gtx.Constraints.Max
	charW := float32(gtx.Sp(unit.Sp(CharFontSize))) * 0.6

	// origin: center for Quick, random position for Hard
	var originX, originY float32
	if g.Mode == ModeHard {
		originX = g.Position.X * float32(rect.X)
		originY = g.Position.Y * float32(rect.Y)
	} else {
		originX = float32(rect.X) / 2
		originY = float32(rect.Y) / 2
	}

	// shake offset
	shakeX := float32(0)
	if g.ShakeTimer > 0 {
		shakeX = float32(math.Sin(float64(g.ShakeTimer)*30.0*2*math.Pi)) * 8.0
	}

	// draw each character in the window
	for offset := -WindowRadius; offset <= WindowRadius; offset++ {
		idx := g.Cursor + offset
		var ch rune = '\u00A0' // non-breaking space for out-of-bounds
		if idx >= 0 && idx < len(g.Text) {
			ch = g.Text[idx]
		}

		absOff := offset
		if absOff < 0 {
			absOff = -absOff
		}
		opacity := float32(1.0)
		switch absOff {
		case 1:
			opacity = 0.75
		case 2:
			opacity = 0.50
		case 3:
			opacity = 0.25
		}
		a := uint8(255 * opacity)

		var col color.NRGBA
		switch {
		case offset == 0:
			col = color.NRGBA{R: 239, G: 68, B: 68, A: a} // red — current char
		case offset > 0:
			col = color.NRGBA{R: 200, G: 200, B: 200, A: a} // gray — upcoming
		case idx >= 0:
			switch g.Correctness[idx] {
			case CorrectnessCorrect:
				col = color.NRGBA{R: 134, G: 239, B: 172, A: a} // green
			case CorrectnessWrong:
				col = color.NRGBA{R: 252, G: 165, B: 165, A: a} // pink
			default:
				col = color.NRGBA{R: 200, G: 200, B: 200, A: a}
			}
		default:
			col = color.NRGBA{R: 200, G: 200, B: 200, A: a}
		}

		x := originX + float32(offset)*charW + shakeX
		y := originY

		drawChar(gtx, th, ch, x, y, CharFontSize, col)

		// underline the current char
		if offset == 0 {
			uy := y + CharFontSize*0.45
			drawLine(gtx, x-charW/2, uy, x+charW/2, uy, col)
		}
	}

	// Hard mode timer
	if g.Mode == ModeHard {
		timerStr := fmt.Sprintf("%.1fs", math.Max(float64(g.RemainingSecs), 0))
		drawTextAt(gtx, th, timerStr, 16, float32(rect.Y)-16, 28, color.NRGBA{R: 239, G: 68, B: 68, A: 255})
	}

	return layout.Dimensions{Size: rect}
}

func drawChar(gtx layout.Context, th *material.Theme, ch rune, x, y, size float32, col color.NRGBA) {
	defer op.Save(gtx.Ops).Load()
	op.Offset(f32.Pt(x, y)).Add(gtx.Ops)
	l := material.Label(th, unit.Sp(size), string(ch))
	l.Color = col
	l.Alignment = text.Middle
	l.Layout(gtx)
}

func drawLine(gtx layout.Context, x1, y1, x2, y2 float32, col color.NRGBA) {
	defer op.Save(gtx.Ops).Load()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	clip.Rect{
		Min: f32.Pt(x1, y1).Round(),
		Max: f32.Pt(x2, y2+2).Round(),
	}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func drawTextAt(gtx layout.Context, th *material.Theme, s string, x, y, size float32, col color.NRGBA) {
	defer op.Save(gtx.Ops).Load()
	op.Offset(f32.Pt(x, y)).Add(gtx.Ops)
	l := material.Label(th, unit.Sp(size), s)
	l.Color = col
	l.Layout(gtx)
}
