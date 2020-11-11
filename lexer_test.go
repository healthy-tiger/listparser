package listparser

import (
	"io"
	"strings"
	"testing"
)

type tokentest struct {
	r    rune
	line int
	col  int
	text string
	err  error
}

var src string = `((abc 123 "hello world")
"hello lf:\nworld cr:\r octet:\073\nhex:\x30\nbackslash:\\\nsingle quote:\'\ndouble quote:\"\nquestion: \?"
(+ 1 2 3) ; comment`

var testresults []tokentest = []tokentest{
	{leftParenthesis, 1, 1, "", nil},
	{leftParenthesis, 1, 2, "", nil},
	{symbol, 1, 3, "abc", nil},
	{' ', 1, 6, "", nil},
	{symbol, 1, 7, "123", nil},
	{' ', 1, 10, "", nil},
	{stringLiteral, 1, 11, "hello world", nil},
	{rightParenthesis, 1, 24, "", nil},
	{stringLiteral, 2, 1, "hello lf:\nworld cr:\r octet:\073\nhex:\x30\nbackslash:\\\nsingle quote:'\ndouble quote:\"\nquestion: ?", nil},
	{leftParenthesis, 3, 1, "", nil},
	{symbol, 3, 2, "+", nil},
	{' ', 3, 3, "", nil},
	{symbol, 3, 4, "1", nil},
	{' ', 3, 5, "", nil},
	{symbol, 3, 6, "2", nil},
	{' ', 3, 7, "", nil},
	{symbol, 3, 8, "3", nil},
	{rightParenthesis, 3, 9, "", nil},
	{' ', 3, 10, "", nil},
	{commentText, 3, 11, " comment", nil},
}

func strcmp(a, b string) (bool, int) {
	if len(a) != len(b) {
		return false, -1
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false, i
		}
	}
	return true, -1
}

func TestToken1(t *testing.T) {
	ss, err := newLexer("TestToken1", strings.NewReader(src))
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < len(testresults); i++ {
		e := testresults[i]
		r, line, col, err := ss.scan()
		if r != e.r || line != e.line || col != e.col || err != e.err {
			t.Error(err)
		}
		switch r {
		case symbol, stringLiteral, commentText:
			st := ss.tokentext()
			if c, index := strcmp(e.text, st); !c {
				t.Errorf("unexpected token text \"%v\", expected \"%v\", different at %d", st, e.text, index)
				t.Errorf("expected char %02x, received %02x", e.text[index], st[index])
			}
		}
	}
	_, _, _, err = ss.scan()
	if err != io.EOF {
		t.Error(err)
	}
}
