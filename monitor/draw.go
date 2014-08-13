package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

var startTime time.Time

func init() {
	startTime = time.Now()
}

func redraw_all() {
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()

	drawCurrentTime(4, 2)

	drawRetCodes(4, 4)
	drawErrStats(30, 4)

	drawPageHits(4, 10)

	midy := h / 2
	midx := (w - edit_box_width) / 2

	// unicode box drawing chars around the edit box
	termbox.SetCell(midx-1, midy, '│', coldef, coldef)
	termbox.SetCell(midx+edit_box_width, midy, '│', coldef, coldef)
	termbox.SetCell(midx-1, midy-1, '┌', coldef, coldef)
	termbox.SetCell(midx-1, midy+1, '└', coldef, coldef)
	termbox.SetCell(midx+edit_box_width, midy-1, '┐', coldef, coldef)
	termbox.SetCell(midx+edit_box_width, midy+1, '┘', coldef, coldef)
	fill(midx, midy-1, edit_box_width, 1, termbox.Cell{Ch: '─'})
	fill(midx, midy+1, edit_box_width, 1, termbox.Cell{Ch: '─'})

	// edit_box.Draw(midx, midy, edit_box_width, 1)
	termbox.SetCursor(midx, midy)

	tbprint(midx+6, midy+3, coldef, coldef, "Press ESC to ʕ◔ϖ◔ʔ")
	termbox.Flush()
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

var knownErrs = []string{
	"parse",
	"nildata",
}

func drawCurrentTime(x, y int) {
	now := time.Now()
	since := now.Sub(startTime)
	h := int(since.Hours())
	m := int(since.Minutes()) % 60
	s := int(since.Seconds()) % 60
	timeStr := fmt.Sprintf("Now:  %-16s    Watching:  %3d:%02d:%02d", now.Format(DATEPRINT), h, m, s)
	for i, c := range timeStr {
		termbox.SetCell(x+i, y, c, coldef, coldef)
	}
}

func drawErrStats(x, y int) {
	colTitle := "Error     Count"
	for i, c := range colTitle {
		termbox.SetCell(x+i, y, c, coldef, coldef)
	}
	y++

	for _, err := range knownErrs {
		errStr := fmt.Sprintf("%-8s  %5d", err, siteStats.ErrStats[err])
		for i, c := range errStr {
			termbox.SetCell(x+i, y, c, coldef, coldef)
		}
		y++
	}
}

var knownCodes = []string{
	"200",
	"404",
}

func drawRetCodes(x, y int) {
	colTitle := "Code      Count"
	for i, c := range colTitle {
		termbox.SetCell(x+i, y, c, coldef, coldef)
	}
	y++

	total := 0
	for _, code := range knownCodes {
		total += siteStats.RetCodes[code]
		errStr := fmt.Sprintf("%-8s  %5d", code, siteStats.RetCodes[code])
		for i, c := range errStr {
			termbox.SetCell(x+i, y, c, coldef, coldef)
		}
		y++
	}

	totalStr := fmt.Sprintf("%-8s  %5d", "total", total)
	for i, c := range totalStr {
		termbox.SetCell(x+i, y, c, coldef, coldef)
	}
}

var knownPages = []string{
	"page1",
	"page2",
	"page3",
	"page4",
}

func drawPageHits(x, y int) {
	colTitle := "Page     Alerts  Count    Hits / min"

	for i, c := range colTitle {
		termbox.SetCell(x+i, y, c, coldef, coldef)
	}
	y++

	hitTotal := 0
	alertTotal := 0
	for _, page := range knownPages {
		hs := siteStats.PageStats[page]
		hist := make([]int, len(hs.HistBins))
		for i := len(hs.HistBins) - 1; i >= 0; i-- {
			hist[i] = hs.HistBins[i].Count
		}

		xcnt := x
		// print page name
		str := fmt.Sprintf("%-8s  ", page)
		for _, c := range str {
			termbox.SetCell(xcnt, y, c, coldef, coldef)
			xcnt++
		}

		// print alerts
		alerts := siteStats.AlertHist[page]
		alertCount := len(alerts)
		acolor := termbox.ColorDefault
		if a := siteStats.OpenAlerts[page]; a != nil {
			acolor = termbox.ColorRed
			alertCount++
		}
		str = fmt.Sprintf("%5d  ", alertCount)
		for _, c := range str {
			termbox.SetCell(xcnt, y, c, coldef, acolor)
			xcnt++
		}

		// print hit infomation
		str = fmt.Sprintf("%5d    [ ", hs.Total)
		for i := len(hist) - 1; i >= 0; i-- {
			str += fmt.Sprintf("%3d ", hist[i])
		}
		str += "]"
		for _, c := range str {
			termbox.SetCell(xcnt, y, c, coldef, coldef)
			xcnt++
		}

		alertTotal += alertCount
		hitTotal += hs.Total
		y++
	}

	totalStr := fmt.Sprintf("%-8s  %5d  %5d", "total", alertTotal, hitTotal)
	for i, c := range totalStr {
		termbox.SetCell(x+i, y, c, coldef, coldef)
	}

}
