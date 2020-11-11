package listparser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// 構文解析のエラーメッセージの定義
const (
	ErrorUnexpectedInputChar            = iota
	ErrorUnexpectedClosingParenthesis   = iota
	ErrorInconsistencyInClosingBrackets = iota
	ErrorTopLevelElementMustBeAList     = iota
	ErrorMissingClosingParenthesis      = iota
	ErrorLexingError                    = iota
)

var errorMessages map[int]string

func init() {
	errorMessages = map[int]string{
		ErrorUnexpectedInputChar:            "Unexpected input char",
		ErrorUnexpectedClosingParenthesis:   "Unexpected closing parenthesis",
		ErrorInconsistencyInClosingBrackets: "Inconsistency in closing brackets",
		ErrorTopLevelElementMustBeAList:     "Top-level element must be a list",
		ErrorMissingClosingParenthesis:      "Missing closing parenthesis",
		ErrorLexingError:                    "Lexing error:",
	}
}

// ParseError パース時のエラーメッセージを格納する
type ParseError struct {
	ErrorLocation Position
	ID            int
	InnerError    error
}

func (err *ParseError) Error() string {
	m := errorMessages[err.ID]
	n := ""
	if err.InnerError != nil {
		n = err.InnerError.Error()
	}
	return fmt.Sprintf("%s:%d:%d %s%s", err.ErrorLocation.Filename, err.ErrorLocation.Line, err.ErrorLocation.Column, m, n)
}

func (err *ParseError) Unwrap() error {
	return err.InnerError
}

func newParseError(filename string, line int, column int, messageid int, innererr error) *ParseError {
	if _, ok := errorMessages[messageid]; !ok {
		panic("Undefined error id")
	}
	return &ParseError{Position{filename, line, column}, messageid, innererr}
}

// Position ソースコード上の位置を表す
type Position struct {
	Filename string
	Line     int
	Column   int
}

// Parse srcをスキャンして*Listの配列を返す。
func Parse(filename string, st *SymbolTable, src io.Reader) ([]*ListElement, error) {
	lists := make([]*ListElement, 0)
	stack := newStack()
	lexer, err := newLexer(filename, src)
	if err != nil {
		return nil, err
	}
	tok, line, column, err := lexer.scan()
	for err == nil {
		toktxt := lexer.tokentext()
		switch tok {
		case symbol:
			lst := stack.peek()
			if lst == nil {
				return nil, newParseError(filename, line, column, ErrorTopLevelElementMustBeAList, nil)
			}
			// IntかFloatとして処理できるか先に確認し、どちらもダメならシンボルにする。
			vi, err := strconv.ParseInt(toktxt, 0, 64)
			if err == nil {
				lst.elements = append(lst.elements, &intElement{vi, Position{filename, line, column}})
			} else {
				vf, err := strconv.ParseFloat(toktxt, 64)
				if err == nil {
					lst.elements = append(lst.elements, &floatElement{vf, Position{filename, line, column}})
				} else {
					lst.elements = append(lst.elements, &symbolIDElement{st.GetSymbolID(toktxt), Position{filename, line, column}})
				}
			}

		case stringLiteral:
			lst := stack.peek()
			if lst == nil {
				return nil, newParseError(filename, line, column, ErrorTopLevelElementMustBeAList, nil)
			}
			lst.elements = append(lst.elements, &stringElement{toktxt, Position{filename, line, column}})

		case commentText:

		default:
			if tok == leftParenthesis || tok == leftSquareBracket || tok == leftCurlyBracket {
				lst := stack.peek()
				lstnew := &ListElement{tok, make([]SyntaxElement, 0), Position{filename, line, column}}
				if lst != nil {
					lst.elements = append(lst.elements, lstnew)
				} else {
					lists = append(lists, lstnew)
				}
				stack.push(lstnew)
			} else if tok == rightParenthesis || tok == rightSquareBracket || tok == rightCurlyBracket {
				lst := stack.peek()
				if lst == nil {
					return nil, newParseError(filename, line, column, ErrorUnexpectedClosingParenthesis, nil)
				} else if !lst.isMatchingParen(tok) {
					return nil, newParseError(filename, line, column, ErrorInconsistencyInClosingBrackets, nil)
				}
				stack.pop()
			} else if tok != tab && tok != space {
				return nil, newParseError(filename, line, column, ErrorUnexpectedInputChar, nil)
			}
		}
		tok, line, column, err = lexer.scan()
	}
	// lexerのエラー＝字句解析のエラーの場合はパースを途中で止める。
	if err != io.EOF {
		return nil, newParseError(filename, line, column, ErrorLexingError, err)
	}
	// スタックが空でないということは閉じていないカッコがあるということ。
	if stack.peek() != nil {
		return nil, newParseError(filename, line, column, ErrorMissingClosingParenthesis, nil)
	}
	return lists, nil
}

// ParseString 文字列をスキャンして*Listの配列を返す。
func ParseString(filename string, st *SymbolTable, src string) ([]*ListElement, error) {
	return Parse(filename, st, strings.NewReader(src))
}

func (p Position) String() string {
	var b bytes.Buffer
	b.WriteString(p.Filename)
	b.WriteString(":")
	b.WriteString(strconv.FormatInt(int64(p.Line), 10))
	b.WriteString(":")
	b.WriteString(strconv.FormatInt(int64(p.Column), 10))
	return b.String()
}
