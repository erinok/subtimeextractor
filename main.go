// Command subtimeextractor reads subtitle times from a file and prints ffmpeg commands to extract the corresponding times.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const slop = 500 * time.Millisecond

func parseTime(x string) time.Duration {
	// 00:28:48,251
	var h, m, s, ms time.Duration
	_, err := fmt.Sscanf(x, "%d:%d:%d,%d", &h, &m, &s, &ms)
	if err != nil {
		fatal("could not parse time:", err)
	}
	return h*time.Hour + m*time.Minute + s*time.Second + ms*time.Millisecond
}

func parseTimes(l string) (time.Duration, time.Duration) {
	// 00:28:48,251 --> 00:28:50,620 ¶ Komm mal rüber.
	l = strings.Replace(l, "-->", "¶", -1)
	parts := strings.Split(l, " ¶ ")
	if len(parts) != 3 {
		fatal("line not in expected format: START --> END ¶ comment:", l)
	}
	return parseTime(parts[0]), parseTime(parts[1])
}

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
		fmt.Println(parseTimes(l))
	}
}

func fatal(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}
