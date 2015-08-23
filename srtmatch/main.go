package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type tm int

const (
	millisecond tm = 1
	second         = 1000 * millisecond
	minute         = 60 * second
	hour           = 60 * minute
)

func parsetm(x string) (tm, error) {
	// 00:28:48,251
	var h, m, s, ms tm
	_, err := fmt.Sscanf(x, "%d:%d:%d,%d", &h, &m, &s, &ms)
	if err != nil {
		return 0, err
	}
	return h*hour + m*minute + s*second + ms*millisecond, nil
}

func parsetmrange(l string) (tm, tm, error) {
	// 00:28:48,251 --> 00:28:50,620
	parts := strings.Split(l, " --> ")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("line not in expected format: START --> END", l)
	}
	if a, err := parsetm(parts[0]); err != nil {
		return 0, 0, err
	} else if b, err := parsetm(parts[1]); err != nil {
		return 0, 0, err
	} else {
		return a, b, nil
	}
}

type xtract struct {
	s string
}

type sub struct {
	a, b tm
	lines []string
}

func parseSRT(fname string) ([]sub, error) {
	var subs []sub
	var cur = &sub{}
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		a, b, err := parsetmrange(l)
		if err == nil {
			if n := len(cur.lines); n > 0 {
				cur.lines = cur.lines[:n - 1]
			}	
			subs = append(subs, sub{a, b, nil})
			cur = &subs[len(subs)-1]
			continue
		}
		cur.lines = append(cur.lines, l)
	}
	return subs, nil
}

func main() {
	subs, err := parseSRT("/Users/erin/de/baader-meinhof-de.srt")
	if err != nil {
		fatal(err)
	}
	for _, sub := range subs {
		fmt.Println(sub.a, sub.b)
		for _, l := range sub.lines {
			fmt.Print("\t", l)
		}
	}
}

func fatal(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}
