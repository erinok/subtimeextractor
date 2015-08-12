// Command subx reads subtitle times from a file and prints a Makefile of ffmpeg commands to extract the corresponding times.
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
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
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

func (t tm) String() string { return formattm(t) }

const clipDir = "/Users/erin/Desktop/declips"
const extractFname = "/Users/erin/de/extract.txt"
const vidFname = "/Users/erin/Downloads/de/baader-meinhof/baader-meinhof.avi"

func getNextFileNum() int {
	files, err := ioutil.ReadDir(clipDir)
	if err != nil {
		fatal(err)
	}
	num := 0
	for _, f := range files {
		var n int
		var ext string
		_, err = fmt.Sscanf(f.Name(), "%d%s", &n, &ext)
		if n > num {
			num = n
		}
	}
	return num+1
}

func main() {
	if len(os.Args) > 1 {
		fatal("usage: subtimeextractor\nextract times from ", extractFname)
	}
	fileNum := getNextFileNum()
	f, err := os.Open(extractFname)
	if err != nil {
		fatal(err)
	}
	r := bufio.NewReader(f)
	fmt.Println("all:\n")
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
		nm := fmt.Sprint("~/Desktop/declips/", fileNum, ".mp3")
		fmt.Print(nm, ":\n")
		fmt.Print("\t", "ffmpeg -y -i ", vidFname, "  -ss ", a, " -to ", b, " ~/Desktop/declips/", fileNum, ".mp3 &> /dev/null\n")
		fmt.Print("all: ", nm, "\n\n")
		fileNum++
	}
}

func fatal(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}
