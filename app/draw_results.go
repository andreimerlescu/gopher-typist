package app

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type NameEntryWidgets struct {
	SubmitBtn widget.Clickable
}

func (w *NameEntryWidgets) Layout(gtx layout.Context, th *material.Theme, r *Router) layout.Dimensions {
	d := r.NameEntry

	if w.SubmitBtn.Clicked(gtx) {
		r.SubmitName()
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(float32(gtx.Constraints.Max.Y) / 5)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				if d.IsNewBest {
					lbl := material.H3(th, "New Best Score!")
					lbl.Color = color.NRGBA{R: 255, G: 215, B: 0, A: 255}
					return lbl.Layout(gtx)
				}
				return material.H3(th, fmt.Sprintf("Level %d Complete!", d.Score.Level+1)).Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				stats := fmt.Sprintf("WPM: %.2f   Accuracy: %.2f%%   CPS: %.2f",
					d.Score.WPM, d.Score.Accuracy, d.Score.CPS)
				return material.Body1(th, stats).Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Body1(th, "Enter your name for the leaderboard:").Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(unit.Dp(240))
				gtx.Constraints.Max.X = gtx.Constraints.Min.X
				ed := material.Editor(th, &nameEditor, "Anonymous")
				return ed.Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(unit.Dp(120))
				gtx.Constraints.Max.X = gtx.Constraints.Min.X
				btn := material.Button(th, &w.SubmitBtn, "Submit")
				btn.Background = color.NRGBA{R: 59, G: 130, B: 246, A: 255}
				return btn.Layout(gtx)
			})
		}),
	)
}
