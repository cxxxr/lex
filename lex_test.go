package lex

import (
	"fmt"
	"strings"
	"testing"
)

const (
	IGNORE = iota
	EOF
	WORD
	NUMBER
	EQUAL
	STRING
)

var lastText string

func test(t *testing.T, f func() int, str string, value int) {
	val := f()
	if val != value {
		t.Errorf("%d != %d", val, value)
	} else if lastText != str {
		t.Errorf("'%s' != '%s'", lastText, str)
	}
}

func makeLexerDef() *LexerDef {
	ld := NewLexerDef()
	ld.SetIgnoreValue(IGNORE)
	ld.SetEOF(func(sc *Scanner) int {
		fmt.Println("<EOF>")
		lastText = ""
		return EOF
	})
	ld.Add("$", func(sc *Scanner) int {
		fmt.Println("<Newline>")
		return IGNORE
	})
	ld.Add("[ \t]+", nil)
	ld.Add("[_a-zA-Z][_a-zA-Z0-9]*", func(sc *Scanner) int {
		lastText = sc.Text()
		fmt.Printf("<WORD:%s>\n", lastText)
		return WORD
	})
	ld.Add("[-+]?[0-9]+", func(sc *Scanner) int {
		lastText = sc.Text()
		fmt.Printf("<NUM:%s>\n", lastText)
		return NUMBER
	})
	ld.Add(`==`, func(sc *Scanner) int {
		lastText = ""
		fmt.Println("<EQUAL>")
		return EQUAL
	})
	ld.Add(`'[^']*'`, func(sc *Scanner) int {
		lastText = sc.Text()
		lastText = lastText[1 : len(lastText)-1]
		fmt.Printf("<STRING:%s>\n", lastText)
		return STRING
	})
	ld.Add(`"[^']*"`, func(sc *Scanner) int {
		lastText = sc.Text()
		lastText = lastText[1 : len(lastText)-1]
		fmt.Printf("<STRING:%s>\n", lastText)
		return STRING
	})
	ld.Add(".", func(sc *Scanner) int {
		lastText = ""
		fmt.Printf("<%c>\n", sc.Text()[0])
		return int(sc.Text()[0])
	})
	return ld
}

func TestLex1(t *testing.T) {
	ld := makeLexerDef()
	f := ld.GenerateLexer(strings.NewReader(`add a, b:
	return a + b
`))
	test(t, f, "add", WORD)
	test(t, f, "a", WORD)
	test(t, f, "", int(','))
	test(t, f, "b", WORD)
	test(t, f, "", int(':'))
	test(t, f, "return", WORD)
	test(t, f, "a", WORD)
	test(t, f, "", int('+'))
	test(t, f, "b", WORD)
	test(t, f, "", EOF)
}

func TestLex2(t *testing.T) {
	ld := makeLexerDef()
	f := ld.GenerateLexer(strings.NewReader(`
fizzbuzz:
	do i = 1 to 100
		if i % 15 == 0 then
			print "fizzbuzz"
		elseif i % 3 == 0 then
			print "fizz"
		elseif i % 5 == 0 then
			print "buzz"
		else
			print i
		end
	end
end
`))

	test(t, f, "fizzbuzz", WORD)
	test(t, f, "", int(':'))
	test(t, f, "do", WORD)
	test(t, f, "i", WORD)
	test(t, f, "", int('='))
	test(t, f, "1", NUMBER)
	test(t, f, "to", WORD)
	test(t, f, "100", NUMBER)
	test(t, f, "if", WORD)
	test(t, f, "i", WORD)
	test(t, f, "", int('%'))
	test(t, f, "15", NUMBER)
	test(t, f, "", EQUAL)
	test(t, f, "0", NUMBER)
	test(t, f, "then", WORD)
	test(t, f, "print", WORD)
	test(t, f, "fizzbuzz", STRING)
	test(t, f, "elseif", WORD)
	test(t, f, "i", WORD)
	test(t, f, "", int('%'))
	test(t, f, "3", NUMBER)
	test(t, f, "", EQUAL)
	test(t, f, "0", NUMBER)
	test(t, f, "then", WORD)
	test(t, f, "print", WORD)
	test(t, f, "fizz", STRING)
	test(t, f, "elseif", WORD)
	test(t, f, "i", WORD)
	test(t, f, "", int('%'))
	test(t, f, "5", NUMBER)
	test(t, f, "", EQUAL)
	test(t, f, "0", NUMBER)
	test(t, f, "then", WORD)
	test(t, f, "print", WORD)
	test(t, f, "buzz", STRING)
	test(t, f, "else", WORD)
	test(t, f, "print", WORD)
	test(t, f, "i", WORD)
	test(t, f, "end", WORD)
	test(t, f, "end", WORD)
	test(t, f, "end", WORD)
	test(t, f, "", EOF)
}
