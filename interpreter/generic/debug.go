package generic

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// SyntaxError is a syntax error
type SyntaxError struct {
	reason      string
	wantedTypes []string
	gotType     string
	line        int
	pos         int
}

func (f *SyntaxError) Error() string {
	return fmt.Sprintf(
		"%v at line %v:%v (wanted=%v, got=%v)",
		f.reason,
		f.line,
		f.pos,
		f.wantedTypes,
		f.gotType,
	)
}

func (f *SyntaxError) String() string { return f.Error() }

// Format implements fmt.Formatter
func (f *SyntaxError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.reason)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.reason)
	case 'q':
		fmt.Fprintf(s, "%q", f.reason)
	}
}

// PrettyPrint a syntax error
func (f *SyntaxError) PrettyPrint(name string, lines []string) string {
	line := f.line
	pos := f.pos

	str := fmt.Sprintln(f.reason)
	str += fmt.Sprintf("In file=%v At=%v:%v\n", name, line, pos)
	str += fmt.Sprintf("Found=%v wanted=%v\n", f.gotType, f.wantedTypes)

	before := []string{}
	about := ""
	after := []string{}
	toRead := len(lines) / 2
	if toRead > 2 {
		before = lines[:toRead-3]
		about = lines[toRead-1]
		after = lines[toRead:]
	} else {
		before = lines[:toRead]
		about = lines[toRead+1]
		after = lines[toRead+2:]
	}

	str += fmt.Sprintln("")
	str += fmt.Sprintln("...")
	line -= toRead - 2 //weird :x
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
	return str
}

// StringSyntaxError is a syntax error in the scope of a Str
type StringSyntaxError struct {
	Filepath string
	Src      string
	SyntaxError
}

// Format implements fmt.Formatter
func (f *StringSyntaxError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('#') {
			io.WriteString(s, f.PrettyPrint())
			return
		}
	}
	f.SyntaxError.Format(s, verb)
}

// PrettyPrint a syntax error
func (f *StringSyntaxError) PrettyPrint() string {
	lines := strings.Split(f.Src, "\n")
	for i := range lines {
		lines[i] += "\n"
	}
	return f.SyntaxError.PrettyPrint(f.Filepath, lines)
}

func (f *StringSyntaxError) Error() string {
	return fmt.Sprintf(
		"in %v %v at line %v:%v (wanted=%v, got=%v)",
		f.Filepath,
		f.reason,
		f.line,
		f.pos,
		f.wantedTypes,
		f.gotType,
	)
}

func (f *StringSyntaxError) String() string { return f.Error() }

// FileSyntaxError is a syntax error in the scope of a File
type FileSyntaxError struct {
	SyntaxError
	Src string
}

// Format implements fmt.Formatter
func (f *FileSyntaxError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('#') {
			io.WriteString(s, f.PrettyPrint())
			return
		}
	}
	f.SyntaxError.Format(s, verb)
}

// PrettyPrint a syntax error
func (f *FileSyntaxError) PrettyPrint() string {
	lines := []string{}
	line := f.line
	from := line - 8
	to := line + 8
	if from < 0 {
		from = 0
	}
	readFileByLine(f.Src, keepLines(from, to, func(line string) {
		lines = append(lines, line)
	}))
	return f.SyntaxError.PrettyPrint(f.Src, lines)
}

func (f *FileSyntaxError) Error() string {
	return fmt.Sprintf(
		"in %v %v at line %v:%v (wanted=%v, got=%v)",
		f.Src,
		f.reason,
		f.line,
		f.pos,
		f.wantedTypes,
		f.gotType,
	)
}

func (f *FileSyntaxError) String() string { return f.Error() }

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
