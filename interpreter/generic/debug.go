package generic

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func DebugSynxtaxError(src string, line, pos int, reason string, got string, wanted ...string) error {
	lines := []string{}
	toRead := 8
	readFileByLine(src, keepLines(line-toRead, line+toRead, func(line string) {
		lines = append(lines, line)
	}))
	before := lines[:toRead-2]
	about := lines[toRead-2]
	after := lines[toRead-1:]
	str := fmt.Sprintln(reason)
	str += fmt.Sprintf("In file=%v At=%v:%v\n", src, line, pos)
	str += fmt.Sprintf("Found=%v wanted=%v\n", got, wanted)
	str += fmt.Sprintln("")
	str += fmt.Sprintln("...")
	line -= toRead - 1
	for _, l := range before {
		line++
		str += fmt.Sprintf("%000d", line)
		str += fmt.Sprint("  ", l)
	}
	line++
	str += fmt.Sprintf("%000d", line)
	str += fmt.Sprint("  ", about)
	str += fmt.Sprintf("   --%v↑\n", strings.Repeat("-", pos))
	for _, l := range after {
		line++
		str += fmt.Sprintf("%000d", line)
		str += fmt.Sprint("  ", l)
	}
	str += fmt.Sprintln("...")
	return errors.New(str)
}

func keepLines(from, to int, h func(line string)) func(line string) {
	c := 0
	return func(line string) {
		if c > from && c < to {
			h(line)
		}
		c++

	}
}

func readFileByLine(filename string, fn func(line string)) error {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return err
	}

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

	var line string
	for {
		line, err = reader.ReadString('\n')
		fn(line)
		if err != nil {
			break
		}
	}
	return err
}
