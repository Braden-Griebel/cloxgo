package compiler

type Scanner struct {
	code    *string
	start   uint
	current uint
	line    uint
}

func initScanner(source *string) *Scanner {
	return &Scanner{code: source, start: 0, current: 0, line: 0}
}
