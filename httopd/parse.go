package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type LineRaw struct {
	line    []byte
	logfile string

	// planned for dealing with multiple &| custom log formats
	// and a generalized parser which may require a more complicated stuct of it's own
	logfmt string // format sting info for parsing (or perhaps the ID of the parser to use)
}

type LineData struct {
	Logfile string // from the logfile `name` in LineRaw.logfile
	Date    time.Time

	RequestStr    string
	RequestMethod string
	SectionStr    string

	Status     string
	ContentLen int

	RemoteHost string
	Rfc931     string
	AuthUser   string
}

const DATELAYOUT = "02/Jan/2006:15:04:05 -0700"
const DATEPRINT = "Jan 02, 2006 15:04:05"

func ParseLineData(raw *LineRaw) (*LineData, error) {
	defer func() {
		if e := recover(); e != nil {
			siteStats.ErrStats["parse"]++
		}
	}()

	var err error
	ld := new(LineData)
	ld.Logfile = raw.logfile
	line := raw.line

	// find date delimiters and parse
	//-------------------------------
	lsb_pos := bytes.IndexRune(line, '[')
	rsb_pos := bytes.IndexRune(line, ']')
	date_str := string(line[lsb_pos+1 : rsb_pos])
	ld.Date, err = time.Parse(DATELAYOUT, date_str)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(line))

	// parse first three fields
	//-------------------------
	last_pos := 0
	curr_pos := bytes.IndexRune(line[last_pos:lsb_pos], ' ') + last_pos
	ld.RemoteHost = string(line[last_pos:curr_pos])
	last_pos = curr_pos + 1

	curr_pos = bytes.IndexRune(line[last_pos:lsb_pos], ' ') + last_pos
	ld.Rfc931 = string(line[last_pos:curr_pos])
	last_pos = curr_pos + 1

	curr_pos = bytes.IndexRune(line[last_pos:lsb_pos], ' ') + last_pos
	ld.AuthUser = string(line[last_pos:curr_pos])

	// parse request string
	//-----------------
	lqt_pos := bytes.IndexRune(line, '"') + 1
	rqt_pos := bytes.IndexRune(line[lqt_pos:], '"') + lqt_pos
	req := line[lqt_pos:rqt_pos]
	ld.RequestStr = string(req)

	// parse RequestMethod from RequestStr
	rrm_pos := bytes.IndexRune(req, ' ')
	ld.RequestMethod = string(req[0:rrm_pos])

	// parse SectionStr from RequestStr
	lreq_pos := bytes.IndexRune(req, '/') + 1
	rreq_pos := bytes.IndexAny(req[lreq_pos:], "/ ") + lreq_pos
	ld.SectionStr = string(req[lreq_pos:rreq_pos])

	// parse status and content-length (the next two fields)
	last_pos = rqt_pos + 2 // skip past '"' and first ' '
	curr_pos = bytes.IndexRune(line[last_pos:], ' ') + last_pos
	ld.Status = string(line[last_pos:curr_pos])
	last_pos = curr_pos + 1
	curr_pos = bytes.IndexRune(line[last_pos:], ' ') + last_pos
	clenStr := string(line[last_pos:curr_pos])
	clenInt, err := strconv.ParseInt(clenStr, 10, 64)
	if err != nil {
		return nil, err
	}
	ld.ContentLen = int(clenInt)

	// there are two remaining fields, which I am omitting

	return ld, nil
}

func startParser(line_chan chan *LineRaw, data_chan chan *LineData) {
	for {
		select {
		case line := <-line_chan:
			// fmt.Print(line)
			ld, err := ParseLineData(line)
			if err != nil {
				fmt.Println(err)
			} else {
				data_chan <- ld
			}
		}
		if quit == true {
			return
		}
	}
}
