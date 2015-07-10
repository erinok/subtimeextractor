// Command subtimeextractor reads subtitle times from a file and prints ffmpeg commands to extract the corresponding times.
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

const slop = 500 * millisecond

func parsetm(x string) tm {
	// 00:28:48,251
	var h, m, s, ms tm
	_, err := fmt.Sscanf(x, "%d:%d:%d,%d", &h, &m, &s, &ms)
	if err != nil {
		fatal("could not parse time:", err)
	}
	return h*hour + m*minute + s*second + ms*millisecond
}

func parsetms(l string) (tm, tm) {
	// 00:28:48,251 --> 00:28:50,620 ¶ Komm mal rüber.
	l = strings.Replace(l, "-->", "¶", -1)
	parts := strings.Split(l, " ¶ ")
	if len(parts) != 3 {
		fatal("line not in expected format: START --> END ¶ comment:", l)
	}
	return parsetm(parts[0]), parsetm(parts[1])
}

func formattm(t tm) string {
	h := t / hour
	m := (t % hour) / minute
	s := t % minute / second
	ms := t % second / millisecond
	return fmt.Sprintf("%02d:%02d:%02d,%03d", h, m, s, ms)
}

func (t tm) String() string { return formattm(t) }

func main() {
	r := bufio.NewReader(os.Stdin)
	if len(os.Args) == 2 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fatal(err)
		}
		r = bufio.NewReader(f)
	}
	for {
		l, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fatal("error reading input:", err)
		}
		a, b := parsetms(l)
		a -= slop
		b += slop
		fmt.Println(a, b)
	}
}

func fatal(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}
