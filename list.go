package listparser

import (
	"fmt"
	"reflect"
)

// SyntaxElement 構文要素を表す。
type SyntaxElement interface {
	Position() Position
	IntValue() (int64, bool)
	FloatValue() (float64, bool)
	StringValue() (string, bool)
	SymbolValue() (SymbolID, bool)
}

// ListElement ListElementまたはValueを0個以上含む
type ListElement struct {
	openchar rune
	elements []SyntaxElement
	pos      Position
}

// SymbolID シンボルのSTreeにおける一意な識別番号
type SymbolID int

const nilInt = 0
const nilFloat = 0.0
const emptyString = ""

func (lst *ListElement) isMatchingParen(close rune) bool {
	if (lst.openchar == leftParenthesis && close == rightParenthesis) ||
		(lst.openchar == leftSquareBracket && close == rightSquareBracket) ||
		(lst.openchar == leftCurlyBracket && close == rightCurlyBracket) {
		return true
	}
	return false
}

// Len lstの子要素の数を返す。
func (lst *ListElement) Len() int {
	return len(lst.elements)
}

// Position lstのソースコード上の位置を返す。
func (lst *ListElement) Position() Position {
	return lst.pos
}

// IntValue lstは整数型の値を持たない。
func (lst *ListElement) IntValue() (int64, bool) {
	return 0, false
}

// FloatValue lstは浮動小数点数型の値を持たない。
func (lst *ListElement) FloatValue() (float64, bool) {
	return 0, false
}

// StringValue lstは文字列型の値を持たない。
func (lst *ListElement) StringValue() (string, bool) {
	return "", false
}

// SymbolValue lstはシンボルではない。
func (lst *ListElement) SymbolValue() (SymbolID, bool) {
	return InvalidSymbolID, false
}

// ElementAt lstのindex番目の要素を返す。
func (lst *ListElement) ElementAt(index int) SyntaxElement {
	if index < 0 || index >= len(lst.elements) {
		return nil
	}
	return lst.elements[index]
}

// IntAt lstのindex番目の要素がint64ならその値を返す。
func (lst *ListElement) IntAt(index int) (int64, bool) {
	se := lst.ElementAt(index)
	if se != nil {
		return se.IntValue()
	}
	return 0, false
}

// FloatAt lstのindex番目の要素がfloat64ならその値を返す。
func (lst *ListElement) FloatAt(index int) (float64, bool) {
	se := lst.ElementAt(index)
	if se != nil {
		return se.FloatValue()
	}
	return 0.0, false
}

// StringAt lstのindex番目の要素がstringならその値を返す。
func (lst *ListElement) StringAt(index int) (string, bool) {
	se := lst.ElementAt(index)
	if se != nil {
		return se.StringValue()
	}
	return "", false
}

// SymbolAt lstのindex番目の要素がint64ならその値を返す。
func (lst *ListElement) SymbolAt(index int) (SymbolID, bool) {
	se := lst.ElementAt(index)
	if se != nil {
		return se.SymbolValue()
	}
	return InvalidSymbolID, false
}

func newLiteral(value interface{}, filename string, line int, column int) SyntaxElement {
	switch v := value.(type) {
	case int64:
		return &intElement{v, Position{filename, line, column}}
	case float64:
		return &floatElement{v, Position{filename, line, column}}
	case SymbolID:
		return &symbolIDElement{v, Position{filename, line, column}}
	case string:
		return &stringElement{v, Position{filename, line, column}}
	}
	panic(fmt.Sprintf("Unexpected value type: %v", reflect.TypeOf(value)))
}

type intElement struct {
	value int64
	pos   Position
}

// Position eのソースコード上の位置を返す。
func (e *intElement) Position() Position {
	return e.pos
}

// IntValue eが整数リテラルなら、整数リテラルのint64型の値を返す。
func (e *intElement) IntValue() (int64, bool) {
	return e.value, true
}

// FloatValue eが浮動小数点数リテラルなら、浮動小数点数リテラルのfloat64の値を返す。
func (e *intElement) FloatValue() (float64, bool) {
	return nilFloat, false
}

// StringValue eが文字列リテラルなら、文字列リテラルのstringの値を返す。
func (e *intElement) StringValue() (string, bool) {
	return emptyString, false
}

// SymbolValue eがシンボルなら、リテラルのSymbolIDを返す。
func (e *intElement) SymbolValue() (SymbolID, bool) {
	return InvalidSymbolID, false
}

type floatElement struct {
	value float64
	pos   Position
}

// Position eのソースコード上の位置を返す。
func (e *floatElement) Position() Position {
	return e.pos
}

// IntValue eが整数リテラルなら、整数リテラルのint64型の値を返す。
func (e *floatElement) IntValue() (int64, bool) {
	return nilInt, false
}

// FloatValue eが浮動小数点数リテラルなら、浮動小数点数リテラルのfloat64の値を返す。
func (e *floatElement) FloatValue() (float64, bool) {
	return e.value, true
}

// StringValue eが文字列リテラルなら、文字列リテラルのstringの値を返す。
func (e *floatElement) StringValue() (string, bool) {
	return emptyString, false
}

// SymbolValue eがシンボルなら、リテラルのSymbolIDを返す。
func (e *floatElement) SymbolValue() (SymbolID, bool) {
	return InvalidSymbolID, false
}

type stringElement struct {
	value string
	pos   Position
}

// Position eのソースコード上の位置を返す。
func (e *stringElement) Position() Position {
	return e.pos
}

// IntValue eが整数リテラルなら、整数リテラルのint64型の値を返す。
func (e *stringElement) IntValue() (int64, bool) {
	return nilInt, false
}

// FloatValue eが浮動小数点数リテラルなら、浮動小数点数リテラルのfloat64の値を返す。
func (e *stringElement) FloatValue() (float64, bool) {
	return nilFloat, false
}

// StringValue eが文字列リテラルなら、文字列リテラルのstringの値を返す。
func (e *stringElement) StringValue() (string, bool) {
	return e.value, true
}

// SymbolValue eがシンボルなら、リテラルのSymbolIDを返す。
func (e *stringElement) SymbolValue() (SymbolID, bool) {
	return InvalidSymbolID, false
}

type symbolIDElement struct {
	value SymbolID
	pos   Position
}

// Position eのソースコード上の位置を返す。
func (e *symbolIDElement) Position() Position {
	return e.pos
}

// IntValue eが整数リテラルなら、整数リテラルのint64型の値を返す。
func (e *symbolIDElement) IntValue() (int64, bool) {
	return nilInt, false
}

// FloatValue eが浮動小数点数リテラルなら、浮動小数点数リテラルのfloat64の値を返す。
func (e *symbolIDElement) FloatValue() (float64, bool) {
	return nilFloat, false
}

// StringValue eが文字列リテラルなら、文字列リテラルのstringの値を返す。
func (e *symbolIDElement) StringValue() (string, bool) {
	return emptyString, false
}

// SymbolValue eがシンボルなら、リテラルのSymbolIDを返す。
func (e *symbolIDElement) SymbolValue() (SymbolID, bool) {
	return e.value, true
}

// IsSymbolID 構文要素eがシンボルID idと等しいかテストする
func IsSymbolID(e SyntaxElement, id SymbolID) bool {
	if sid, ok := e.SymbolValue(); ok && sid == id {
		return true
	}
	return false
}

// IsSymbol 構文要素eがシンボルかどうかテストする
func IsSymbol(e SyntaxElement) bool {
	_, ok := e.SymbolValue()
	return ok
}

// IsList 構文要素eがリストかどうかテストする
func IsList(e SyntaxElement) bool {
	_, ok := e.(*ListElement)
	return ok
}

// IsInt 構文要素eが整数かどうかテストする
func IsInt(e SyntaxElement) bool {
	_, ok := e.(*intElement)
	return ok
}

// IsFloat 構文要素eが浮動小数点数かどうかテストする
func IsFloat(e SyntaxElement) bool {
	_, ok := e.(*floatElement)
	return ok
}

// IsString 構文要素eが文字列かどうかテストする
func IsString(e SyntaxElement) bool {
	_, ok := e.(*stringElement)
	return ok
}
