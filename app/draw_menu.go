package app

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type MenuWidgets struct {
	QuickBtn widget.Clickable
	HardBtn  widget.Clickable
}

func (m *MenuWidgets) Layout(gtx layout.Context, th *material.Theme, r *Router) layout.Dimensions {
	if m.QuickBtn.Clicked(gtx) {
		r.StartGame(ModeQuick)
	}
	if m.HardBtn.Clicked(gtx) {
		r.StartGame(ModeHard)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(float32(gtx.Constraints.Max.Y) / 4)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H2(th, "Gopher Typist").Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				sub := material.Body2(th, "press Escape anytime to return here")
				sub.Color = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
				return sub.Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(160))
						gtx.Constraints.Max.X = gtx.Constraints.Min.X
						btn := material.Button(th, &m.QuickBtn, "New Quick Game")
						btn.Background = color.NRGBA{R: 59, G: 130, B: 246, A: 255}
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(160))
						gtx.Constraints.Max.X = gtx.Constraints.Min.X
						btn := material.Button(th, &m.HardBtn, "New Hard Game")
						btn.Background = color.NRGBA{R: 239, G: 68, B: 68, A: 255}
						return btn.Layout(gtx)
					}),
				)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(40)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return DrawLeaderboard(gtx, th, r)
			})
		}),
	)
}
