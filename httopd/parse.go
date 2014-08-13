package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type LineData struct {
	RemoteHost string
	Date       time.Time
	RequestStr string
	SectionStr string
	Status     string
	ContentLen int

	Rfc931   string
	AuthUser string
}

const DATELAYOUT = "02/Jan/2006:15:04:05 -0700"
const DATEPRINT = "Jan 02, 2006 15:04:05"

func ParseLineData(line []byte) (*LineData, error) {
	var err error
	ld := new(LineData)

	defer func() {
		if e := recover(); e != nil {
			siteStats.ErrStats["parse"]++
		}
	}()

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
	// fmt.Printf("  %d %d : %q  %q\n", last_pos, curr_pos, string(line[0:lsb_pos]), string(line[last_pos:curr_pos]))
	ld.RemoteHost = string(line[last_pos:curr_pos])
	last_pos = curr_pos + 1

	curr_pos = bytes.IndexRune(line[last_pos:lsb_pos], ' ') + last_pos
	// fmt.Printf("  %d %d : %q  %q\n", last_pos, curr_pos, string(line[last_pos:lsb_pos]), string(line[last_pos:curr_pos]))
	ld.Rfc931 = string(line[last_pos:curr_pos])
	last_pos = curr_pos + 1

	curr_pos = bytes.IndexRune(line[last_pos:lsb_pos], ' ') + last_pos
	// fmt.Printf("  %d %d : %q  %q\n", last_pos, curr_pos, string(line[last_pos:lsb_pos]), string(line[last_pos:curr_pos]))
	ld.AuthUser = string(line[last_pos:curr_pos])

	// parse request string
	//-----------------
	lqt_pos := bytes.IndexRune(line, '"') + 1
	rqt_pos := bytes.IndexRune(line[lqt_pos:], '"') + lqt_pos
	ld.RequestStr = string(line[lqt_pos:rqt_pos])
	// parse SectionStr from RequestStr
	req := line[lqt_pos:rqt_pos]
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

func startParser(line_chan chan []byte, data_chan chan *LineData) {
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
