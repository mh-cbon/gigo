package generic

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// SyntaxError is a syntax error
type SyntaxError struct {
	reason error
	line   int
	pos    int
}

// NewSyntaxError creates a new syntax error of reason r and pos l:p
func NewSyntaxError(r error, l, p int) SyntaxError {
	return SyntaxError{
		reason: r,
		line:   l,
		pos:    p,
	}
}

func (f *SyntaxError) Error() string {
	return fmt.Sprintf(
		"%v at line %v:%v",
		f.reason,
		f.line,
		f.pos,
	)
}

func (f *SyntaxError) String() string { return f.Error() }

// Format implements fmt.Formatter
func (f *SyntaxError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			m := fmt.Sprintf(
				"%v at line %v:%v\n\n%+v",
				f.reason,
				f.line,
				f.pos,
				f.reason,
			)
			io.WriteString(s, m)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.Error())
	case 'q':
		fmt.Fprintf(s, "%q", f.Error())
	}
}

// PrettyPrint a syntax error
func (f *SyntaxError) PrettyPrint(lines []string) string {
	line := f.line
	pos := f.pos

	str := ""
	// str := fmt.Sprintln(f.reason)
	// str += fmt.Sprintf("In file=%v At=%v:%v\n", name, line, pos)
	// str += fmt.Sprintf("Found=%v wanted=%v\n", f.gotType, f.wantedTypes)

	toRead := len(lines) / 2
	before := lines[:toRead]
	about := lines[toRead]
	after := lines[toRead+1:]
	line -= toRead
	if toRead > 2 {
		before = lines[0 : toRead-1]
		about = lines[toRead-1]
		after = lines[toRead:]
	} else {
		line -= 2
	}

	str += fmt.Sprintln("")
	if line > 0 {
		str += fmt.Sprintln("...")
	}
	for _, l := range before {
		str += fmt.Sprintf("%000d", line)
		str += fmt.Sprint("  ", l)
		line++
	}
	str += fmt.Sprintf("%000d", line)
	str += fmt.Sprint("  ", about)
	line++
	if pos < 0 {
		str += fmt.Sprintf("✘%v", strings.Repeat("-", len(strconv.Itoa(line))-1))
		str += fmt.Sprintf("- ↑↑↑ ???\n")
	} else {
		str += fmt.Sprintf("✘%v", strings.Repeat("-", len(strconv.Itoa(line))-1))
		str += fmt.Sprintf("%v↑\n", strings.Repeat("-", pos))
	}
	for _, l := range after {
		str += fmt.Sprintf("%000d", line)
		str += fmt.Sprint("  ", l)
		line++
	}
	str += fmt.Sprintln("...")
	return str
}

// ParseError is an error about parsing
type ParseError struct {
	SyntaxError
	wantedTypes []string
	gotType     string
}

//NewParseError creates a parse error
func NewParseError(reason error, n Tokener, got string, wanted []string) *ParseError {
	return &ParseError{
		SyntaxError: SyntaxError{
			reason: reason,
			line:   n.GetPos().Line,
			pos:    n.GetPos().Pos,
		},
		wantedTypes: wanted,
		gotType:     got,
	}
}

func (f *ParseError) Error() string {
	return fmt.Sprintf(
		"%v (wanted=%v, got=%v)",
		f.SyntaxError.Error(),
		f.wantedTypes,
		f.gotType,
	)
}

// PrettyPrint a syntax error
func (f *ParseError) PrettyPrint(lines []string) string {
	str := fmt.Sprintln(f.Error())
	str += fmt.Sprintf("\n\n%v", f.SyntaxError.PrettyPrint(lines))
	return str
}

// StringSyntaxError is a syntax error in the scope of a Str
type StringSyntaxError struct {
	ParseError
	Filepath string
	Src      string
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
	str := fmt.Sprintln(f.Error())
	str += fmt.Sprintf("\n\n%v", f.ParseError.PrettyPrint(lines))
	return str
}

func (f *StringSyntaxError) Error() string {
	return fmt.Sprintf("in %v %v", f.Filepath, f.ParseError.Error())
}

func (f *StringSyntaxError) String() string { return f.Error() }

// FileSyntaxError is a syntax error in the scope of a File
type FileSyntaxError struct {
	ParseError
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
	from := line - 8 // something weird here.
	to := line + 8
	if from < 0 {
		from = 0
	}
	readFileByLine(f.Src, keepLines(from, to, func(line string) {
		lines = append(lines, line)
	}))
	str := fmt.Sprintln(f.Error())
	str += fmt.Sprintf("\n\n%v", f.ParseError.PrettyPrint(lines))
	return str
}

func (f *FileSyntaxError) Error() string {
	return fmt.Sprintf("in %v %v", f.Src, f.ParseError.Error())
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

var tplErrLine = regexp.MustCompile(`:([0-9]+):`)

//NewStringTplSyntaxError creates a new syntax error for a template
func NewStringTplSyntaxError(from error, name, tplContent string) *StringTplSyntaxError {
	msg := from.Error()
	res := tplErrLine.FindAllStringSubmatch(msg, -1)
	line := 0
	if len(res) > 0 {
		if x, err := strconv.Atoi(res[0][1]); err == nil {
			line = x
		}
	}
	return &StringTplSyntaxError{
		SyntaxError: NewSyntaxError(from, line, -1),
		Name:        name,
		Src:         tplContent,
	}
}

// StringTplSyntaxError is a syntax error in the scope of a Str
type StringTplSyntaxError struct {
	SyntaxError
	Name string
	Src  string
}

// Format implements fmt.Formatter
func (f *StringTplSyntaxError) Format(s fmt.State, verb rune) {
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
func (f *StringTplSyntaxError) PrettyPrint() string {
	lines := strings.Split(f.Src, "\n")
	for i := range lines {
		lines[i] += "\n"
	}
	from := f.line - 4
	to := f.line + 4
	if to > len(lines) {
		to = len(lines)
	}
	if from < 0 {
		from = 0
	}
	lines = lines[from:to]
	str := fmt.Sprintln(f.Error())
	str += fmt.Sprintf("\n\n%v", f.SyntaxError.PrettyPrint(lines))
	return str
}

func (f *StringTplSyntaxError) Error() string {
	return fmt.Sprintf("in %v %v", f.Name, f.SyntaxError.Error())
}

func (f *StringTplSyntaxError) String() string { return f.Error() }
