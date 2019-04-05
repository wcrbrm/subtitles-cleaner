package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// Timeframe contains beginning and end of the frame
type Timeframe struct {
	StartTime string
	EndTime   string
}

// ToTime can convert string into seconds
func ToTime(tm string) (int, error) {
	// cut "," and what is after it
	commaIndex := strings.Index(tm, ",")
	if commaIndex > -1 {
		tm = strings.TrimSpace(tm[:commaIndex])
	}
	// split using ":" as delimiter
	s := strings.Split(tm, ":")

	h, _ := strconv.Atoi(s[0])
	min, _ := strconv.Atoi(s[1])
	sec, _ := strconv.Atoi(s[2])
	return h*60*60 + min*60 + sec, nil
}

// SecondsBetween can calculate Seconds between
func (tf *Timeframe) SecondsBetween(tf2 *Timeframe) int {
	t1, err1 := ToTime(tf.EndTime)
	t2, err2 := ToTime(tf2.StartTime)
	if err1 != nil || err2 != nil {
		return 0
	}
	return t2 - t1
}

// NewTimeFrame constructs time frame from line
func NewTimeFrame(s string) *Timeframe {
	sepIndex := strings.Index(s, "-->")
	if sepIndex == -1 {
		return nil
	}
	after := strings.TrimSpace(s[sepIndex+3:])
	before := strings.TrimSpace(s[:sepIndex])
	return &Timeframe{before, after}
}

// IsTimeFrame returns if this is a valid time frame string
func IsTimeFrame(s string) bool {
	sepIndex := strings.Index(s, "-->")
	return strings.Count(s, ":") == 4 && sepIndex > -1
}

// Row represents each subtitle
type Row struct {
	Num   int
	TF    *Timeframe
	Lines []string
}

func parseFile(filename string, fileout string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fOut, err := os.Create(fileout)
	if err != nil {
		log.Fatal(err)
	}
	defer fOut.Close()
	w := bufio.NewWriter(fOut)

	rows := make([]Row, 0)
	var r *Row
	scanner := bufio.NewScanner(file)
	// file scanner goes through every line of input.txt
	ln := 0
	for scanner.Scan() {
		ln++
		line := strings.TrimSpace(scanner.Text())
		if r == nil {
			// we were expecting a number to start
			num, errNum := strconv.Atoi(line)
			if errNum != nil {
				fmt.Println("ERROR: Number expected, line #" + strconv.Itoa(ln))
			} else {
				r = &Row{num, nil, make([]string, 0)}
			}
		} else if r.TF == nil {
			if !IsTimeFrame(line) {
				fmt.Println("ERROR: TimeFrame expected, line #" + strconv.Itoa(ln))
				r = nil
			} else {
				r.TF = NewTimeFrame(line)
			}
		} else if line != "" {
			r.Lines = append(r.Lines, line)
		} else { // if empty line
			if r != nil {
				rows = append(rows, *r) // copied
			}
			r = nil
		}
	}

	for index, ro := range rows {

		secs := 0
		if index > 0 {
			secs = rows[index-1].TF.SecondsBetween(ro.TF)
		}
		withSpace := index > 0 && secs > 10
		if withSpace {
			// fmt.Println("")
			w.WriteString("\n")
		}
		// fmt.Println("#", ro.Num, "secs=", secs) // comment it
		for _, l := range ro.Lines {
			// fmt.Println(l)
			w.WriteString(l + "\n")
		}
	}
	w.Flush()
}

func main() {
	dir := "./s01"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.Contains(f.Name(), ".srt") {
			outFile := strings.Replace(f.Name(), ".srt", ".txt", -1)
			fmt.Println(f.Name(), "-->", outFile)
			parseFile(dir+"/"+f.Name(), dir+"/"+outFile)
		}
	}
}
