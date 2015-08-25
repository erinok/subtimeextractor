package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
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

func formattm(t tm) string {
	h := t / hour
	m := (t % hour) / minute
	s := t % minute / second
	ms := t % second / millisecond
	return fmt.Sprintf("%02d:%02d:%02d,%03d", h, m, s, ms)
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
	a, b  tm
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
				cur.lines = cur.lines[:n-1]
			}
			subs = append(subs, sub{a, b, nil})
			cur = &subs[len(subs)-1]
			continue
		}
		cur.lines = append(cur.lines, l)
	}
	return subs, nil
}

var spaceRE = regexp.MustCompile(`[[:space:]]+`)
var punctRE = regexp.MustCompile(`[^\pL ]+`)

func sanitize(s string) string {
	// collapse space so space fixes in german don't mess us up; a smarter edit-distance match would be a better fix
	s = spaceRE.ReplaceAllString(s, "")
	s = punctRE.ReplaceAllString(s, "")
	s = strings.ToLower(s)
	return s
}

func collapseLines(lines []string) string {
	return sanitize(strings.Join(lines, " "))
}

type attr struct {
	idx  int
	tm   tm
	open bool
}

type astring struct {
	text  string
	attrs []attr
}

func (a *astring) append(s string, start, stop tm) {
	a.attrs = append(a.attrs, attr{len(a.text), start, true})
	a.text += s
	a.attrs = append(a.attrs, attr{len(a.text), stop, false})
}

func (a *astring) lookupRange(i, j int) (start, stop attr) {
	for _, attr := range a.attrs {
		if attr.open {
			if attr.idx <= i {
				start = attr
			}
		} else {
			if j <= attr.idx {
				stop = attr
				return
			}
		}
	}
	return
}

func matchscore(a, b string) int {
	if len(a) != len(b) {
		panic(fmt.Sprintln("bad args", a, b))
	}
	s := 0
	for i := range b {
		if a[i] != b[i] {
			s++
		}
	}
	return s
}

// find the index i s.t. matchscore(a, txt[i:len(a)]) is minimal
func bestmatch(a, txt string) int {
	best_s := len(txt)
	best_i := 0
	for i := 0; i <= len(txt)-len(a); i++ {
		if s := matchscore(a, txt[i:i+len(a)]); s < best_s {
			best_s = s
			best_i = i
		}
	}
	return best_i
}

func main() {
	subs, err := parseSRT("/Users/erin/de/baader-meinhof-de.srt")
	if err != nil {
		fatal(err)
	}
	var text astring
	for _, sub := range subs {
		text.append(collapseLines(sub.lines), sub.a, sub.b)
	}
	for _, raw := range readlines("/Users/erin/de/lines.txt") {
		line := sanitize(raw)
		k := bestmatch(line, text.text)
		start, stop := text.lookupRange(k, k+len(line))
		fmt.Print(formattm(start.tm), " --> ", formattm(stop.tm), " ", raw)
	}
}

func fatal(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func readlines(fname string) (lines []string) {
	f, err := os.Open(fname)
	if err != nil {
		fatal(err)
	}
	r := bufio.NewReader(f)
	for {
		l, err := r.ReadString('\n')
		if err == io.EOF {
			return lines
		} else if err != nil {
			fatal(err)
		}
		lines = append(lines, l)
	}
	return
}
