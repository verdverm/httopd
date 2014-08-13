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

var w, h int

func redraw_all() {
	termbox.Clear(coldef, coldef)
	w, h = termbox.Size()

	drawCurrentTime(1, 0)

	drawRetCodes(1, 2)
	drawErrStats(30, 2)

	if alertDetailsView {
		drawAlertDetails(1, 7)
	} else {
		drawPageStats(1, 7)
	}

	drawFooter()

	// midy := h / 2
	// midx := (w - edit_box_width) / 2

	termbox.HideCursor()

	tbprint(w-6, h-1, coldef, termbox.ColorBlue, "ʕ◔ϖ◔ʔ")
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
	timeStr := fmt.Sprintf("Now:  %-16s  Watching:  %3d:%02d:%02d", now.Format(DATEPRINT), h, m, s)
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

var selectedRow = 0
var maxSelectedRow = len(knownPages) - 1

func drawPageStats(x, y int) {
	colTitle := "Page     Alerts  Count    Hits / min"

	for i := 0; i < w; i++ {
		termbox.SetCell(i, y, ' ', coldef, termbox.ColorCyan)
	}
	for i, c := range colTitle {
		termbox.SetCell(x+i, y, c, coldef, termbox.ColorCyan)
	}
	y++

	if selectedRow < 0 {
		selectedRow = 0
	}

	for p, page := range knownPages {
		hs := siteStats.PageStats[page]
		hist := make([]int, len(hs.HistBins))
		for i := len(hs.HistBins) - 1; i >= 0; i-- {
			hist[i] = hs.HistBins[i].Count
		}

		xcnt := x

		fg_col, bg_col := coldef, coldef
		if p == selectedRow {
			fg_col = termbox.ColorBlack
			bg_col = termbox.ColorYellow
		}

		// print page name
		str := fmt.Sprintf("%-8s  ", page)
		for i := 0; i < w; i++ {
			termbox.SetCell(i, y, ' ', fg_col, bg_col)
		}
		for _, c := range str {
			termbox.SetCell(xcnt, y, c, fg_col, bg_col)
			xcnt++
		}

		// print alerts
		alerts := siteStats.AlertHist[page]
		alertCount := len(alerts)
		alert_fg := fg_col
		alert_bg := bg_col
		if a := siteStats.OpenAlerts[page]; a != nil {
			alert_fg = termbox.ColorDefault
			alert_bg = termbox.ColorRed
			alertCount++
		}
		str = fmt.Sprintf("%5d  ", alertCount)
		for _, c := range str {
			termbox.SetCell(xcnt, y, c, alert_fg, alert_bg)
			xcnt++
		}

		// print hit infomation
		str = fmt.Sprintf("%5d    [ ", hs.Total)
		for i := len(hist) - 1; i >= 0; i-- {
			str += fmt.Sprintf("%3d ", hist[i])
		}
		str += "]"
		for _, c := range str {
			termbox.SetCell(xcnt, y, c, fg_col, bg_col)
			xcnt++
		}

		y++
	}

}

const alertDateFormat = "01-02 15:04"

func drawAlertDetails(x, y int) {
	colTitle := "Type           Start          End            Details"
	for i := 0; i < w; i++ {
		termbox.SetCell(i, y, ' ', coldef, termbox.ColorCyan)
	}
	for i, c := range colTitle {
		termbox.SetCell(x+i, y, c, coldef, termbox.ColorCyan)
	}
	y++

	if selectedRow < 0 {
		selectedRow = 0
	}

	fg_col, bg_col := coldef, coldef

	if alert := siteStats.OpenAlerts[alertPage]; alert != nil {
		bstr, estr := alert.BeginTime.Format(alertDateFormat), "  "
		if alert.BeginTime.Before(alert.EndTime) {
			estr = alert.EndTime.Format(alertDateFormat)
		}
		str := fmt.Sprintf("%-12s   %-12s   %-12s   %s", alert.Type, bstr, estr, alert.Detail)
		for i, c := range str {
			termbox.SetCell(x+i, y, c, fg_col, bg_col)
		}
		y++
	}

	alerts := siteStats.AlertHist[alertPage]
	for P := len(alerts) - 1; P >= 0; P-- {
		alert := alerts[P]
		hs := siteStats.PageStats[alertPage]
		hist := make([]int, len(hs.HistBins))
		for i := len(hs.HistBins) - 1; i >= 0; i-- {
			hist[i] = hs.HistBins[i].Count
		}

		fg_col, bg_col := coldef, coldef

		// print alert
		str := fmt.Sprintf("%-12s   %-12s   %-12s   %s", alert.Type, alert.BeginTime.Format(alertDateFormat), alert.EndTime.Format(alertDateFormat), alert.Detail)
		for i, c := range str {
			termbox.SetCell(x+i, y, c, fg_col, bg_col)
		}
		y++
	}

}

func drawFooter() {
	footerText := " Esc:Back   Ctrl-Q:Quit   Enter:Detail" //"<_sort_> "
	for i := 0; i < w; i++ {
		termbox.SetCell(i, h-1, ' ', coldef, termbox.ColorCyan)
	}
	for i, c := range footerText {
		col := termbox.ColorCyan
		if c != ' ' {
			col = termbox.ColorBlue
		}
		termbox.SetCell(i, h-1, c, coldef, col)
	}

}
