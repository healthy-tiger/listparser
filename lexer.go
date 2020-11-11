package listparser

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode/utf8"
)

// 字句解析エラー
var (
	ErrorIllegalLexerState        = errors.New("Illegal lexer state")
	ErrorUnexpectedEndOfLine      = errors.New("Unexpected end of line")
	ErrorIllegalCharacterEncoding = errors.New("Illegal character encoding")
	ErrorIllegalEscapeSequence    = errors.New("Illegal escape sequence '%c'")
)

const (
	ctxString    = iota
	ctxEscSeq    = iota
	ctxEscOctet1 = iota
	ctxEscOctet2 = iota
	ctxEscHex    = iota
	ctxEscHex1   = iota
)

const (
	backslash   = '\\'
	doublequote = '"'
	semicolon   = ';'
)

const (
	tokTab                = '\t'
	tokSpace              = ' '
	tokLeftParenthesis    = '('
	tokRightParenthesis   = ')'
	tokLeftSquareBracket  = '['
	tokRightSquareBracket = ']'
	tokLeftCurlyBracket   = '{'
	tokRightCurlyBracket  = '}'
)

var stdEscSeq = map[rune]rune{
	'a':  '\x07',
	'b':  '\x08',
	'f':  '\x0c',
	'n':  '\x0a',
	'r':  '\x0d',
	't':  '\x09',
	'v':  '\x0b',
	'\\': '\\',
	'\'': '\'',
	'"':  '"',
	'?':  '?',
}

var octValues = map[rune]int32{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7,
}

var hexValues = map[rune]int32{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15,
	'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15,
}

type slexer struct {
	inputname   string
	linescanner *bufio.Scanner
	reader      io.RuneScanner
	lasttext    string
	line        int
	column      int
}

func (ss *slexer) nextline() error {
	if ss.linescanner.Scan() {
		ss.reader = strings.NewReader(ss.linescanner.Text())
		ss.line = ss.line + 1
		ss.column = 1
		return nil
	}
	e := ss.linescanner.Err()
	if e == nil {
		e = io.EOF
	}
	return e
}

// readString 文字列リテラルの最初の'"'以降の部分をエスケープシーケンスを解釈して文字列を返す。
func (ss *slexer) readString() (string, int, error) {
	runes := make([]rune, 0)
	stat := ctxString
	nr := 0
	var oct int32
	var hex int32
	r, sz, err := ss.reader.ReadRune()
	for r != utf8.RuneError && err == nil {
		nr++
		switch stat {
		case ctxString:
			if r == backslash {
				stat = ctxEscSeq
			} else if r == doublequote {
				return string(runes), nr, nil
			} else {
				runes = append(runes, r)
			}

		case ctxEscSeq:
			ec, ok := stdEscSeq[r]
			if ok {
				stat = ctxString
				runes = append(runes, ec)
			} else {
				ov, ok := octValues[r]
				if ok {
					stat = ctxEscOctet1
					oct = ov
				} else if r == 'x' {
					stat = ctxEscHex
				} else {
					return "", nr, ErrorIllegalEscapeSequence
				}
			}

		case ctxEscOctet1:
			ov, ok := octValues[r]
			if ok {
				stat = ctxEscOctet2
				oct = oct*8 + ov
			} else {
				return "", nr, ErrorIllegalEscapeSequence
			}

		case ctxEscOctet2:
			ov, ok := octValues[r]
			if ok {
				stat = ctxString
				oct = oct*8 + ov
				runes = append(runes, oct)
			} else {
				return "", nr, ErrorIllegalEscapeSequence
			}

		case ctxEscHex:
			hv, ok := hexValues[r]
			if ok {
				stat = ctxEscHex1
				hex = hv
			} else {
				return "", nr, ErrorIllegalEscapeSequence
			}

		case ctxEscHex1:
			hv, ok := hexValues[r]
			if ok {
				stat = ctxString
				hex = hex*16 + hv
				runes = append(runes, hex)
			} else {
				return "", nr, ErrorIllegalEscapeSequence
			}
		}
		r, sz, err = ss.reader.ReadRune()
	}
	if sz == 0 || err == io.EOF {
		return "", nr, ErrorUnexpectedEndOfLine
	}
	// 行単位で処理しているので行末以外は符号化のエラーとしていいはず。
	return "", nr, ErrorIllegalCharacterEncoding
}

func (ss *slexer) readSymbol() (string, int, error) {
	rs := make([]rune, 0)
	nr := 0
	r, sz, err := ss.reader.ReadRune()
	for r != utf8.RuneError && err == nil {
		nr++
		switch r {
		case tokTab, tokSpace, semicolon, tokLeftParenthesis, tokLeftSquareBracket, tokLeftCurlyBracket, tokRightParenthesis, tokRightSquareBracket, tokRightCurlyBracket:
			// 空白かコメントかカッコ（開く又は閉じる）まで読む。
			// 読み込んじゃった一文字はUnreadRune()で戻しておく
			err = ss.reader.UnreadRune()
			if err != nil {
				return "", nr, err
			}
			nr--
			return string(rs), nr, nil

		default:
			rs = append(rs, r)
		}
		r, sz, err = ss.reader.ReadRune()
	}

	if sz == 0 || err == io.EOF { // 行の末尾まで読み込んだ場合、読み込んだ部分までをシンボルとして返す。
		return string(rs), nr, nil
	}
	// 行単位で処理しているので行末以外は符号化のエラーとしていいはず。
	return "", nr, ErrorIllegalCharacterEncoding
}

func (ss *slexer) readComment() (string, int, error) {
	// 行末まで読み込んで返す。
	rs := make([]rune, 0)
	nr := 0
	r, sz, err := ss.reader.ReadRune()
	for r != utf8.RuneError && err == nil {
		nr++
		rs = append(rs, r)
		r, sz, err = ss.reader.ReadRune()
	}

	if sz == 0 || err == io.EOF { // readerの末尾に達した場合
		return string(rs), nr, nil
	}
	// 行単位で処理しているので行末以外は符号化のエラーとしていいはず。
	return "", nr, ErrorIllegalCharacterEncoding
}

func newLexer(inputname string, reader io.Reader) (*slexer, error) {
	ss := &slexer{inputname, bufio.NewScanner(reader), nil, "", 0, 0}
	err := ss.nextline()
	if err != nil {
		return nil, err
	}
	return ss, nil
}

const (
	symbol        = -(iota + 1)
	stringLiteral = -(iota + 1)
	commentText   = -(iota + 1)
)

// scan 次のトークンを読み込む
// 読み込んだ文字またはトークンの種類、行番号、列番号、エラー（ある場合は）を返す。
func (ss *slexer) scan() (rune, int, int, error) {
	// 現在の行の次の一文字を読み込む
	r, sz, err := ss.reader.ReadRune()

	// 行の末尾まで読み込んだ場合は次の空じゃない行までスキップする。
	for sz == 0 || err == io.EOF {
		err = ss.nextline()
		if err != nil {
			return 0, ss.line, ss.column, err
		}
		r, sz, err = ss.reader.ReadRune()
	}
	if r == utf8.RuneError && sz > 0 {
		return 0, ss.line, ss.column, ErrorIllegalCharacterEncoding
	} else if err != nil {
		// EOF以外のエラーの場合、行単位で処理しているのでないはず。
		return 0, ss.line, ss.column, ErrorIllegalCharacterEncoding
	}

	switch r {
	case tokLeftParenthesis, tokLeftSquareBracket, tokLeftCurlyBracket, tokRightParenthesis, tokRightSquareBracket, tokRightCurlyBracket, tokTab, tokSpace:
		c := ss.column
		ss.column = ss.column + 1
		return r, ss.line, c, nil
	case doublequote:
		sl, nr, err := ss.readString()
		c := ss.column
		ss.column = ss.column + 1 + nr // '"'の分はss.readString()の返り値には含まれないので+1
		if err == nil {
			ss.lasttext = sl
			return stringLiteral, ss.line, c, nil
		}
		return 0, ss.line, c, err
	case semicolon:
		cm, _, err := ss.readComment()
		if err == nil {
			ss.lasttext = cm
			return commentText, ss.line, ss.column, nil
		}
		return 0, ss.line, ss.column, err
	default:
		err = ss.reader.UnreadRune()
		if err != nil {
			return 0, ss.line, ss.column, ErrorIllegalLexerState
		}
		sl, nr, err := ss.readSymbol()
		c := ss.column
		ss.column = ss.column + nr
		if err == nil {
			ss.lasttext = sl
			return symbol, ss.line, c, nil
		}
		return 0, ss.line, c, err
	}
}

func (ss *slexer) tokentext() string {
	return ss.lasttext
}
