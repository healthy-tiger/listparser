package listparser

import (
	"errors"
	"fmt"
)

// 内部エラーの定義
var (
	ErrorArgumentIsNil            = errors.New("Argument is nil")
	ErrorValueTypeIsNotAsExpected = errors.New("Value type is not as expected")
	ErrorUndefinedSymbol          = errors.New("Undefined symbol")
)

// 構文解析、字句解析のエラーメッセージの定義
const (
	ErrorUnmatchedParenthesis           = iota
	ErrorUnexpectedToken                = iota
	ErrorUnexpectedInputChar            = iota
	ErrorInsufficientInput              = iota
	ErrorFirstElementTypeMustBeASymbol  = iota
	ErrorStringLiteralMustBeASingleLine = iota
	ErrorIllegalEscapeSequence          = iota
	ErrorNotStringLiteral               = iota
	ErrorTopLevelElementMustBeAList     = iota
	ErrorMissingClosingParenthesis      = iota
)

var errorMessages map[int]string

func init() {
	errorMessages = map[int]string{
		ErrorUnmatchedParenthesis:           "Unmatched parenthesis",
		ErrorUnexpectedToken:                "Unexpected token",
		ErrorUnexpectedInputChar:            "Unexpected input char '%c'",
		ErrorInsufficientInput:              "Insufficient input",
		ErrorFirstElementTypeMustBeASymbol:  "First element type must be a symbol",
		ErrorStringLiteralMustBeASingleLine: "String literal must be a single line",
		ErrorIllegalEscapeSequence:          "Illegal escape sequence '%c'",
		ErrorNotStringLiteral:               "Not string literal",
		ErrorTopLevelElementMustBeAList:     "Top-level element must be a list",
		ErrorMissingClosingParenthesis:      "Missing closing parenthesis",
	}
}

// ParseError パース時のエラーメッセージを格納する
type ParseError struct {
	ErrorLocation Position
	ID            int
	Arg           interface{}
}

func (err *ParseError) Error() string {
	h := fmt.Sprintf("%s:%d:%d ", err.ErrorLocation.Filename, err.ErrorLocation.Line, err.ErrorLocation.Column)
	m := errorMessages[err.ID]
	if err.Arg != nil {
		m = fmt.Sprintf(m, err.Arg)
	}
	return h + m
}

func newError(filename string, line int, column int, messageid int, arg interface{}) *ParseError {
	if _, ok := errorMessages[messageid]; !ok {
		panic("Undefined error id")
	}
	return &ParseError{Position{filename, line, column}, messageid, arg}
}
