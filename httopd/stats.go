package main

import (
	"fmt"
	"time"
)

type Recorder func(*HistStats, *LineData)
type Trigger func(*HistStats) bool

type SiteStats struct {
	// use PageStats here to avoid duplication of history data stuct
	ErrStats   map[string]int
	RetCodes   map[string]int
	PageStats  map[string]HistStats
	PageErrors map[string]HistStats

	OpenAlerts map[string]*PageAlert
	AlertHist  map[string][]*PageAlert
}

var (
	BinSize    = time.Minute // size of a bin
	HistoryLen = 10          // number of bins
)

type HistBin struct {
	Start time.Time
	Count int
}

type HistStats struct {
	HistBins []HistBin
	LastTime time.Time
	Total    int
}

type PageAlert struct {
	Type   string
	Detail string

	BeginTime time.Time
	EndTime   time.Time
}

var siteStats SiteStats

func init() {
	siteStats.ErrStats = make(map[string]int)
	siteStats.RetCodes = make(map[string]int)

	siteStats.PageStats = make(map[string]HistStats)
	siteStats.AlertHist = make(map[string][]*PageAlert)
	siteStats.OpenAlerts = make(map[string]*PageAlert)
}

const ALERT_THRESHOLD = 105

func startStats(data_chan chan *LineData) {

	for {
		select {
		case ld := <-data_chan:
			if ld == nil {
				siteStats.ErrStats["nildata"]++
				continue
			}

			updateStats(ld)
			checkAlerts(ld)

		}
	}
}

func updateStats(ld *LineData) {
	code := ld.Status
	page := ld.SectionStr
	siteStats.RetCodes[code]++

	hs := siteStats.PageStats[page]
	hs.Total++
	if hs.LastTime.Minute() != ld.Date.Minute() {
		// time for new bin
		if len(hs.HistBins) > 9 {
			hs.HistBins = hs.HistBins[1:]
		}
		d := ld.Date
		ldTime := time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Minute(), 0, time.UTC)
		bin := HistBin{
			Start: ldTime,
			Count: 1,
		}
		hs.HistBins = append(hs.HistBins, bin)
		hs.LastTime = ldTime
	} else {
		// continue with current bin
		l := len(hs.HistBins) - 1
		hs.HistBins[l].Count++
	}
	siteStats.PageStats[page] = hs
}

func checkAlerts(ld *LineData) {
	page := ld.SectionStr
	hs := siteStats.PageStats[page]
	bins := hs.HistBins
	l := len(bins)

	// make sure we have enough history
	if l < 3 {
		return
	}

	// calc trigger for last two COMPLETE minutes
	ave := (bins[l-3].Count + bins[l-2].Count) / 2

	alert := siteStats.OpenAlerts[page]
	if alert == nil && ave >= ALERT_THRESHOLD {
		// check to open a new alert
		alert = new(PageAlert)
		alert.Type = "High Traffic"
		alert.Detail = fmt.Sprintf("hits per minute > %d", ALERT_THRESHOLD)
		d := ld.Date
		ldTime := time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Minute(), 0, time.UTC)
		alert.BeginTime = ldTime
		siteStats.OpenAlerts[page] = alert
	} else if alert != nil && ave <= ALERT_THRESHOLD &&
		alert.BeginTime.Minute() != ld.Date.Minute() {
		// check to close the current alert
		d := ld.Date
		ldTime := time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Minute(), 0, time.UTC)
		alert.EndTime = ldTime
		alerts := siteStats.AlertHist[page]
		alerts = append(alerts, alert)
		siteStats.AlertHist[page] = alerts
		siteStats.OpenAlerts[page] = nil
	}

}
