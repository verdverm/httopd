package main

import (
	"fmt"
	"time"
)

const ALERT_THRESHOLD = 35

type Recorder func(*HistStats, *LineData)
type Trigger func(*HistStats) bool

// main "DS"
type MainDS struct {
	LogNames []string
	Logs     map[string]*SiteStats
	ErrStats map[string]int

	// this alerts is not in use, but is planned to move all alerts here
	AlertHist map[string][]*PageAlert
}

var siteStats MainDS

func init() {
	siteStats.Logs = make(map[string]*SiteStats)
	siteStats.ErrStats = make(map[string]int)
}

type SiteStats struct {
	LogName string
	// use PageStats here to avoid duplication of history data stuct
	RetCodes   map[string]int
	PageStats  map[string]HistStats
	PageErrors map[string]HistStats

	OpenAlerts map[string]*PageAlert
	AlertHist  map[string][]*PageAlert
}

// called in startWatcher
func addSiteStats(logfn string) {
	ss := new(SiteStats)
	ss.LogName = logfn

	// not thread safe
	siteStats.LogNames = append(siteStats.LogNames, logfn)
	siteStats.Logs[logfn] = ss

	ss.RetCodes = make(map[string]int)

	ss.PageStats = make(map[string]HistStats)
	ss.AlertHist = make(map[string][]*PageAlert)
	ss.OpenAlerts = make(map[string]*PageAlert)
}

var (
	BinSize    = time.Minute // size of a bin
	HistoryLen = 10          // number of bins
)

type HistBin struct {
	Start time.Time
	Bytes int
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
	logfn := ld.Logfile
	siteStats.Logs[logfn].RetCodes[code]++

	hs := siteStats.Logs[logfn].PageStats[page]
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
		hs.HistBins[l].Bytes += ld.ContentLen
	}
	siteStats.Logs[logfn].PageStats[page] = hs
}

func checkAlerts(ld *LineData) {
	page := ld.SectionStr
	logfn := ld.Logfile

	hs := siteStats.Logs[logfn].PageStats[page]
	bins := hs.HistBins
	l := len(bins)

	// make sure we have enough history
	if l < 3 {
		return
	}

	// calc trigger for last two COMPLETE minutes
	ave := (bins[l-3].Count + bins[l-2].Count) / 2

	alert := siteStats.Logs[logfn].OpenAlerts[page]
	if alert == nil && ave >= ALERT_THRESHOLD {
		// check to open a new alert
		alert = new(PageAlert)
		alert.Type = "High Traffic"
		alert.Detail = fmt.Sprintf("hits per minute > %d", ALERT_THRESHOLD)
		d := ld.Date
		ldTime := time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Minute(), 0, time.UTC)
		alert.BeginTime = ldTime
		siteStats.Logs[logfn].OpenAlerts[page] = alert
	} else if alert != nil && ave <= ALERT_THRESHOLD &&
		alert.BeginTime.Minute() != ld.Date.Minute() {
		// check to close the current alert
		d := ld.Date
		ldTime := time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Minute(), 0, time.UTC)
		alert.EndTime = ldTime
		alerts := siteStats.Logs[logfn].AlertHist[page]
		alerts = append(alerts, alert)
		siteStats.Logs[logfn].AlertHist[page] = alerts
		siteStats.Logs[logfn].OpenAlerts[page] = nil
	}

}
