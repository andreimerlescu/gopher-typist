package app

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func DrawTimeout(gtx layout.Context, th *material.Theme, r *Router, w *Widgets) layout.Dimensions {
	if w.RetryBtn.Clicked(gtx) {
		r.StartGame(r.Mode)
	}
	if w.TimeoutMenu.Clicked(gtx) {
		r.GoToMenu()
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(float32(gtx.Constraints.Max.Y) / 3)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H1(th, "Out of time!").Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H5(th, "Type faster!").Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(130))
						gtx.Constraints.Max.X = gtx.Constraints.Min.X
						btn := material.Button(th, &w.RetryBtn, "Retry")
						btn.Background = color.NRGBA{R: 239, G: 68, B: 68, A: 255}
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(130))
						gtx.Constraints.Max.X = gtx.Constraints.Min.X
						btn := material.Button(th, &w.TimeoutMenu, "Main Menu")
						btn.Background = color.NRGBA{R: 75, G: 75, B: 75, A: 255}
						return btn.Layout(gtx)
					}),
				)
			})
		}),
	)
}
