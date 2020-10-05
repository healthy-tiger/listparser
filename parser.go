package listparser

import (
	"bytes"
	"io"
	"strconv"
	"strings"
)

// MinSymbolID シンボルIDの最小値
const MinSymbolID = 10

// InvalidSymbolID 無効なシンボルID(-1)
const InvalidSymbolID = MinSymbolID - MinSymbolID - 1

// SymbolTable シンボルIDとシンボル名のマップ
type SymbolTable struct {
	symbolMap map[string]SymbolID
}

// NewSymbolTable 新しいSymbolTableを作る。
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{make(map[string]SymbolID)}
}

// GetSymbolID はシンボルnameに対するIDを返す。
// IDが割り当てられていないシンボルに対しては、新たにIDを割り当てて返す。
func (st *SymbolTable) GetSymbolID(name string) SymbolID {
	n, ok := st.symbolMap[name]
	if !ok {
		n = SymbolID(len(st.symbolMap) + MinSymbolID)
		st.symbolMap[name] = SymbolID(n)
	}
	return n
}

// GetSymbolName はシンボルのIDからシンボル名を取得する。
func (st *SymbolTable) GetSymbolName(id SymbolID) (string, error) {
	for k, v := range st.symbolMap {
		if v == id {
			return k, nil
		}
	}
	return "", ErrorUndefinedSymbol
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
	tokenizer, err := newTokenizer(filename, src)
	if err != nil {
		return nil, err
	}
	tok, line, column, err := tokenizer.scan()
	for err == nil {
		toktxt := tokenizer.tokentext()
		switch tok {
		case symbol:
			lst := stack.peek()
			if lst == nil {
				return nil, newError(filename, line, column, ErrorTopLevelElementMustBeAList, nil)
			}
			// IntかFloatとして処理できるか先に確認し、どちらもダメならシンボルにする。
			vi, err := strconv.ParseInt(toktxt, 0, 64)
			if err == nil {
				lst.elements = append(lst.elements, newLiteral(vi, filename, line, column))
			} else {
				vf, err := strconv.ParseFloat(toktxt, 64)
				if err == nil {
					lst.elements = append(lst.elements, newLiteral(vf, filename, line, column))
				} else {
					lst.elements = append(lst.elements, newLiteral(st.GetSymbolID(toktxt), filename, line, column))
				}
			}

		case stringLiteral:
			lst := stack.peek()
			if lst == nil {
				return nil, newError(filename, line, column, ErrorTopLevelElementMustBeAList, nil)
			}
			lst.elements = append(lst.elements, newLiteral(toktxt, filename, line, column))

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
				if lst == nil || !lst.isMatchingParen(tok) {
					return nil, newError(filename, line, column, ErrorUnexpectedInputChar, tok)
				}
				stack.pop()
			} else if tok != tab && tok != space {
				return nil, newError(filename, line, column, ErrorUnexpectedInputChar, tok)
			}
		}
		tok, line, column, err = tokenizer.scan()
	}
	// tokenizerのエラー＝字句解析のエラーの場合はパースを途中で止める。
	if err != io.EOF {
		return nil, err
	}
	if stack.peek() != nil {
		return nil, newError(filename, line, column, ErrorMissingClosingParenthesis, nil)
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
