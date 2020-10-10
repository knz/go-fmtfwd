package fmtfwd

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ReproducePrintf formats the value of arg using the current formatting
// parameters in the provided fmt.State and verb, into the provided
// io.Writer.
//
// When implementing a Format() method, one can use
//
//     ReproducePrintf(s, s, verb, ...)
//
// Where the fmt.State is used both as the source of formatting
// parameters and the output destination.
func ReproducePrintf(w io.Writer, s fmt.State, verb rune, arg interface{}) {
	justV, revFmt := MakeFormat(s, verb)
	if justV {
		// Common case, avoids generating then parsing the format again.
		fmt.Fprint(w, arg)
	} else {
		fmt.Fprintf(w, revFmt, arg)
	}
}

// MakeFormat is a helper to aid with the implementation of
// fmt.Formatter for custom types. It reproduces the format currently
// active in fmt.State and verb. This is provided because Go's
// standard fmt.State does not make the original format string
// available to us.
//
// If the return value justV is true, then the current state
// was found to be %v exactly; in that case the caller
// can avoid a full-blown Printf call and use just Print instead
// to take a shortcut.
func MakeFormat(s fmt.State, verb rune) (justV bool, format string) {
	plus, minus, hash, sp, z := s.Flag('+'), s.Flag('-'), s.Flag('#'), s.Flag(' '), s.Flag('0')
	w, wp := s.Width()
	p, pp := s.Precision()

	if !plus && !minus && !hash && !sp && !z && !wp && !pp {
		switch verb {
		case 'v':
			return true, "%v"
		case 's':
			return false, "%s"
		case 'd':
			return false, "%d"
		}
		// Other cases handled in the slow path below.
	}

	var f strings.Builder
	f.WriteByte('%')
	if plus {
		f.WriteByte('+')
	}
	if minus {
		f.WriteByte('-')
	}
	if hash {
		f.WriteByte('#')
	}
	if sp {
		f.WriteByte(' ')
	}
	if z {
		f.WriteByte('0')
	}
	if wp {
		f.WriteString(strconv.Itoa(w))
	}
	if pp {
		f.WriteByte('.')
		f.WriteString(strconv.Itoa(p))
	}
	f.WriteRune(verb)
	return false, f.String()
}
