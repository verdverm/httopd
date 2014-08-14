package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

var startTime time.Time

func init() {
	startTime = time.Now()
}

var w, h int

var tmpLogFn = "/home/tony/gocode/src/github.com/verdverm/httopd/logs/host.httopd-1.access.log"

var knownPages = []string{
	"page1",
	"page2",
	"page3",
	"page4",
}

var selectedRow = 0
var colHeaderRow = 7
var minSelectedRow = colHeaderRow + 1
var maxSelectedRow = colHeaderRow + 1

func redraw_all() {
	termbox.Clear(coldef, coldef)
	w, h = termbox.Size()

	drawCurrentTime(1, 0)

	// temporary, this is changing
	drawRetCodes(1, 2)
	drawErrStats(30, 2)

	// convert to 'detailsView'
	if alertDetailsView {
		// determine details to draw
		alertPage = "page1"
		ss := siteStats.Logs[tmpLogFn]

		drawSectionDetails(1, 7, alertPage, ss)
		// draw page details (more detailed stats & alert hist)
		// or
		// draw logfile details (aggregates of the page details)
	} else {
		drawColumnHeaders(1, colHeaderRow)
		y := colHeaderRow + 1
		for _, f := range siteStats.LogNames {
			ss := siteStats.Logs[f]
			drawPageStats(1, y, ss)
			y += len(ss.PageStats)
		}
		maxSelectedRow = y - 1
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
	timeStr := fmt.Sprintf("Now:  %-24s  Watching:  %3d:%02d:%02d", now.Format(DATEPRINT), h, m, s)
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
	selRow := fmt.Sprintf("%-8s  %5d  %5d  %5d", "selRow", selectedRow, minSelectedRow, maxSelectedRow)
	for i, c := range selRow {
		termbox.SetCell(x+i, y, c, coldef, coldef)
	}

}

var knownCodes = []string{
	"200",
	"404",
}

func drawRetCodes(x, y int) {
	// temporary, want to calc global here
	ss := siteStats.Logs[tmpLogFn]

	colTitle := "Code      Count"
	for i, c := range colTitle {
		termbox.SetCell(x+i, y, c, coldef, coldef)
	}
	y++

	total := 0
	for _, code := range knownCodes {
		total += ss.RetCodes[code]
		errStr := fmt.Sprintf("%-8s  %5d", code, ss.RetCodes[code])
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

func drawColumnHeaders(x, y int) {
	columnHeaders := fmt.Sprintf(
		"%-4s  %-24s  %-6s   %-6s   %-48s",
		"CID", "Page", "Alerts", "Count", "Hits / min",
	)

	for i := 0; i < w; i++ {
		termbox.SetCell(i, y, ' ', coldef, termbox.ColorBlue)
	}
	for i, c := range columnHeaders {
		termbox.SetCell(x+i, y, c, coldef, termbox.ColorBlue)
	}
}
func drawPageStats(x, y int, ss *SiteStats) {

	if selectedRow < minSelectedRow {
		selectedRow = minSelectedRow
	}

	// draw log row

	fg_col, bg_col := coldef, coldef
	if y == selectedRow {
		fg_col = termbox.ColorBlack
		bg_col = termbox.ColorYellow
	}

	// print log file name
	lpos := strings.LastIndex(ss.LogName, "/") + 1
	lfn_short := ss.LogName[lpos:]
	str := fmt.Sprintf("%-4d  %-24s  ", y, lfn_short)
	for i := 0; i < w; i++ {
		termbox.SetCell(i, y, ' ', fg_col, bg_col)
	}
	xcnt := x
	for _, c := range str {
		termbox.SetCell(xcnt, y, c, fg_col, bg_col)
		xcnt++
	}
	y++

	// draw page sections
	for _, page := range knownPages {
		hs := ss.PageStats[page]
		hist := make([]int, len(hs.HistBins))
		for i := len(hs.HistBins) - 1; i >= 0; i-- {
			hist[i] = hs.HistBins[i].Count
		}

		xcnt := x

		fg_col, bg_col = coldef, coldef
		if y == selectedRow {
			fg_col = termbox.ColorBlack
			bg_col = termbox.ColorYellow
		}

		// print page name
		str := fmt.Sprintf("%-4d  %-24s  ", y, "  "+page)
		for i := 0; i < w; i++ {
			termbox.SetCell(i, y, ' ', fg_col, bg_col)
		}
		for _, c := range str {
			termbox.SetCell(xcnt, y, c, fg_col, bg_col)
			xcnt++
		}

		// print alerts
		alerts := ss.AlertHist[page]
		alertCount := len(alerts)
		alert_fg := fg_col
		alert_bg := bg_col
		if a := ss.OpenAlerts[page]; a != nil {
			alert_fg = termbox.ColorDefault
			alert_bg = termbox.ColorRed
			alertCount++
		}
		str = fmt.Sprintf("%6d  ", alertCount)
		for _, c := range str {
			termbox.SetCell(xcnt, y, c, alert_fg, alert_bg)
			xcnt++
		}

		// print hit infomation
		str = fmt.Sprintf("%6d    [ ", hs.Total)
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

func drawSectionDetails(x, y int, alertsPage string, ss *SiteStats) {

	// draw the details
	colTitle := "Type           Start          End            Details"
	for i := 0; i < w; i++ {
		termbox.SetCell(i, y, ' ', coldef, termbox.ColorBlue)
	}
	for i, c := range colTitle {
		termbox.SetCell(x+i, y, c, coldef, termbox.ColorBlue)
	}
	y++

	if selectedRow < 0 {
		selectedRow = 0
	}

	fg_col, bg_col := coldef, coldef

	if alert := ss.OpenAlerts[alertPage]; alert != nil {
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

	alerts := ss.AlertHist[alertPage]
	for P := len(alerts) - 1; P >= 0; P-- {
		alert := alerts[P]
		hs := ss.PageStats[alertPage]
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
		termbox.SetCell(i, h-1, ' ', coldef, termbox.ColorBlue)
	}
	for i, c := range footerText {
		termbox.SetCell(i, h-1, c, coldef, termbox.ColorBlue)
	}

}
