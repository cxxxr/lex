package lex

import (
	"bufio"
	"io"
	"regexp"
)

type Action func(*Scanner) int

type LexerDef struct {
	regs        []*regexp.Regexp
	actions     []Action
	eofAction   Action
	ignoreValue int
}

func NewLexerDef() *LexerDef {
	lexer := new(LexerDef)
	lexer.regs = make([]*regexp.Regexp, 0)
	lexer.actions = make([]Action, 0)
	return lexer
}

func (ld *LexerDef) SetEOF(action Action) {
	ld.eofAction = action
}

func (ld *LexerDef) SetIgnoreValue(value int) {
	ld.ignoreValue = value
}

func (ld *LexerDef) Add(pattern string, action Action) {
	ld.regs = append(ld.regs, regexp.MustCompile(pattern))
	ld.actions = append(ld.actions, action)
}

type Scanner struct {
	scanner    *bufio.Scanner
	line       string
	pos        int
	ahead      string
	matchedLen int
}

func (scanner *Scanner) update() bool {
	scanner.matchedLen = 0
	if len(scanner.line) < scanner.pos {
		if scanner.scanner.Scan() {
			scanner.line = scanner.scanner.Text()
			scanner.pos = 0
		} else {
			return false
		}
	}
	scanner.ahead = scanner.line[scanner.pos:]
	return true
}

func (scanner *Scanner) step(n int) {
	scanner.pos += n
	scanner.matchedLen = n
}

func (scanner *Scanner) scan(r *regexp.Regexp) int {
	indexes := r.FindStringIndex(scanner.ahead)
	if len(indexes) != 0 && indexes[0] == 0 {
		return indexes[1]
	} else {
		return -1
	}
}

func findPattern(ld *LexerDef, scanner *Scanner) (Action, int, bool) {
	var action Action
	max := -1
	for i := 0; i < len(ld.regs); i++ {
		n := scanner.scan(ld.regs[i])
		if n != -1 && max < n {
			if n == 0 {
				n = 1
			}
			max = n
			action = ld.actions[i]
		}
	}
	return action, max, max != -1
}

func (scanner *Scanner) Text() string {
	return scanner.ahead[:scanner.matchedLen]
}

func (ld *LexerDef) GenerateLexer(r io.Reader) func() int {
	scanner := &Scanner{
		scanner: bufio.NewScanner(r),
		line:    "",
		pos:     1,
	}
	return func() int {
		for {
			if !scanner.update() {
				return ld.eofAction(scanner)
			}
			action, n, found := findPattern(ld, scanner)
			if found {
				scanner.step(n)
				if action != nil {
					val := action(scanner)
					if val != ld.ignoreValue {
						return val
					}
				}
			} else {
				scanner.step(1)
			}
		}
	}
}
