package listparser

import (
	"fmt"
	"reflect"
)

// ElementKind 構文要素の種類を表す型
type ElementKind int

// 構文要素の種類を表す。シンボルIDの下限値よりも下の整数にしているのでシンボルIDと混ぜて使うことができる。
const (
	Symbol ElementKind = MinSymbolID - 1
	Int    ElementKind = MinSymbolID - 2
	Float  ElementKind = MinSymbolID - 3
	String ElementKind = MinSymbolID - 4
	List   ElementKind = MinSymbolID - 5
)

// SyntaxElement 構文要素を表す。
type SyntaxElement interface {
	Position() Position
	IsList() bool
	IntValue() (int64, bool)
	FloatValue() (float64, bool)
	StringValue() (string, bool)
	SymbolValue() (SymbolID, bool)
	ElementAt(int) SyntaxElement
	Kind() ElementKind
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

// IsList lstがリストの場合はtrueを返す。
func (lst *ListElement) IsList() bool {
	return true
}

// Kind 要素の種類を返す。
func (lst *ListElement) Kind() ElementKind {
	return List
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

// Matches 子要素の種類またはシンボルIDが引数patに合致するかテストする。
func (lst *ListElement) Matches(pat ...interface{}) bool {
	if lst.Len() != len(pat) {
		return false
	}
	for i, p := range pat {
		switch p.(type) {
		case ElementKind:
			if lst.ElementAt(i).Kind() != p.(ElementKind) {
				return false
			}
		case SymbolID:
			if s, ok := lst.ElementAt(i).SymbolValue(); ok {
				if s != p.(SymbolID) {
					return false
				}
			} else {
				return false
			}
		}
	}
	return true
}

// StartWith 子要素の種類またはシンボルIDが引数patで始まるかテストする。
func (lst *ListElement) StartWith(pat ...interface{}) bool {
	if lst.Len() < len(pat) {
		return false
	}
	for i, p := range pat {
		switch p.(type) {
		case ElementKind:
			if lst.ElementAt(i).Kind() != p.(ElementKind) {
				return false
			}
		case SymbolID:
			if s, ok := lst.ElementAt(i).SymbolValue(); ok {
				if s != p.(SymbolID) {
					return false
				}
			} else {
				return false
			}
		}
	}
	return true
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

// IsList eがリストならtrueを返す。
func (e *intElement) IsList() bool {
	return false
}

// Position eのソースコード上の位置を返す。
func (e *intElement) Position() Position {
	return e.pos
}

func (e *intElement) Kind() ElementKind {
	return Int
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

func (e *intElement) ElementAt(_ int) SyntaxElement {
	return nil
}

type floatElement struct {
	value float64
	pos   Position
}

// IsList eがリストならtrueを返す。
func (e *floatElement) IsList() bool {
	return false
}

// Position eのソースコード上の位置を返す。
func (e *floatElement) Position() Position {
	return e.pos
}

func (e *floatElement) Kind() ElementKind {
	return Float
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

func (e *floatElement) ElementAt(_ int) SyntaxElement {
	return nil
}

type stringElement struct {
	value string
	pos   Position
}

// IsList eがリストならtrueを返す。
func (e *stringElement) IsList() bool {
	return false
}

// Position eのソースコード上の位置を返す。
func (e *stringElement) Position() Position {
	return e.pos
}

func (e *stringElement) Kind() ElementKind {
	return String
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

func (e *stringElement) ElementAt(_ int) SyntaxElement {
	return nil
}

type symbolIDElement struct {
	value SymbolID
	pos   Position
}

// IsList eがリストならtrueを返す。
func (e *symbolIDElement) IsList() bool {
	return false
}

// Position eのソースコード上の位置を返す。
func (e *symbolIDElement) Position() Position {
	return e.pos
}

func (e *symbolIDElement) Kind() ElementKind {
	return Symbol
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

func (e *symbolIDElement) ElementAt(_ int) SyntaxElement {
	return nil
}
