package app

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func DrawNameEntry(gtx layout.Context, th *material.Theme, r *Router, w *Widgets) layout.Dimensions {
	if w.SubmitBtn.Clicked(gtx) {
		r.NameInput = w.NameEditor.Text()
		r.SubmitName()
		w.NameEditor.SetText("")
	}

	d := r.NameEntry

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
				return material.Editor(th, &w.NameEditor, "Anonymous").Layout(gtx)
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

func DrawScore(gtx layout.Context, th *material.Theme, r *Router, w *Widgets) layout.Dimensions {
	if w.NextLevelBtn.Clicked(gtx) {
		r.NextLevel()
	}
	if w.PlayAgainBtn.Clicked(gtx) {
		r.StartGame(r.Mode)
	}
	if w.MenuBtn.Clicked(gtx) {
		r.GoToMenu()
	}

	d := r.Score

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(float32(gtx.Constraints.Max.Y) / 5)}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.H3(th, fmt.Sprintf("You ranked #%d on the leaderboard!", d.PlayerRank)).Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return DrawLeaderboard(gtx, th, r)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					func() layout.FlexChild {
						if r.HasNextLevel() {
							return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min.X = gtx.Dp(unit.Dp(150))
								gtx.Constraints.Max.X = gtx.Constraints.Min.X
								btn := material.Button(th, &w.NextLevelBtn, "Next Level")
								btn.Background = color.NRGBA{R: 59, G: 130, B: 246, A: 255}
								return btn.Layout(gtx)
							})
						}
						return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{}
						})
					}(),
					layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(150))
						gtx.Constraints.Max.X = gtx.Constraints.Min.X
						btn := material.Button(th, &w.PlayAgainBtn, "Play Again")
						btn.Background = color.NRGBA{R: 34, G: 197, B: 94, A: 255}
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(150))
						gtx.Constraints.Max.X = gtx.Constraints.Min.X
						btn := material.Button(th, &w.MenuBtn, "Main Menu")
						btn.Background = color.NRGBA{R: 75, G: 75, B: 75, A: 255}
						return btn.Layout(gtx)
					}),
				)
			})
		}),
	)
}
