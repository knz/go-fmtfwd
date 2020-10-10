# go-fmtfwd

Perfect forwarding for Go printf formatting strings

This library facilitates the implementation of `fmt.Formatter` on
custom types, when there is a need to forward the call
to a default value.

For example:

```go
func (t *T) Format(s fmt.State, verb rune) {
  if s.Flag('#') {
     fmt.Fprint(s, "hello")
  } else {
     _, fmtForward := fmtfwd.MakeFormat(s, verb)
     fmt.Fprintf(s, fmtForward, "world")
  }
}
```

Alternatively, the combination `MakeFormat` + `fmt.Fprintf` call
can be combined as a single call `fmtfwd.ReproducePrintf(s, s, verb,
...)`.

Package documentation here: https://pkg.go.dev/github.com/knz/go-fmtfwd
