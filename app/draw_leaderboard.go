package app

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// column proportions must sum to 1.0
var colProportions = [7]float32{0.08, 0.18, 0.10, 0.13, 0.10, 0.08, 0.33}
var colHeaders = [7]string{"Rank", "Name", "WPM", "Accuracy", "CPS", "Level", "Date (UTC)"}

func DrawLeaderboard(gtx layout.Context, th *material.Theme, r *Router) layout.Dimensions {
	if len(r.Leaderboard.Entries) == 0 {
		lbl := material.Body1(th, "No scores yet — play a game!")
		lbl.Color = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
		return lbl.Layout(gtx)
	}

	tableW := float32(gtx.Constraints.Max.X)
	if tableW > 700 {
		tableW = 700
	}

	colW := [7]float32{}
	for i, p := range colProportions {
		colW[i] = p * tableW
	}

	rowH := unit.Dp(30)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.H5(th, "Leaderboard").Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),

		// header row
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				headerCells(gtx, th, colW)...,
			)
		}),

		// separator
		layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),

		// data rows
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				dataRows(gtx, th, r, colW, rowH)...,
			)
		}),
	)
}

func headerCells(gtx layout.Context, th *material.Theme, colW [7]float32) []layout.FlexChild {
	children := make([]layout.FlexChild, len(colHeaders))
	for i, h := range colHeaders {
		w := colW[i]
		header := h
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = int(w)
			gtx.Constraints.Max.X = int(w)
			lbl := material.Body2(th, header)
			lbl.Color = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
			return lbl.Layout(gtx)
		})
	}
	return children
}

func dataRows(gtx layout.Context, th *material.Theme, r *Router, colW [7]float32, rowH unit.Dp) []layout.FlexChild {
	children := make([]layout.FlexChild, len(r.Leaderboard.Entries))
	for i, entry := range r.Leaderboard.Entries {
		e := entry
		rank := i + 1
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			rowColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			if rank == 1 {
				rowColor = color.NRGBA{R: 255, G: 215, B: 0, A: 255}
			}
			cells := []string{
				fmt.Sprintf("#%d", rank),
				e.Name,
				fmt.Sprintf("%.1f", e.WPM),
				fmt.Sprintf("%.1f%%", e.Accuracy),
				fmt.Sprintf("%.2f", e.CPS),
				fmt.Sprintf("%d", e.Level+1),
				e.TimestampDisplay(),
			}
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				rowCells(gtx, th, cells, colW, rowH, rowColor)...,
			)
		})
	}
	return children
}

func rowCells(gtx layout.Context, th *material.Theme, cells []string, colW [7]float32, rowH unit.Dp, col color.NRGBA) []layout.FlexChild {
	children := make([]layout.FlexChild, len(cells))
	for i, cell := range cells {
		w := colW[i]
		text := cell
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = int(w)
			gtx.Constraints.Max.X = int(w)
			gtx.Constraints.Min.Y = gtx.Dp(rowH)
			lbl := material.Body2(th, text)
			lbl.Color = col
			return lbl.Layout(gtx)
		})
	}
	return children
}
