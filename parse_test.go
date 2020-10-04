package listparser

import "testing"

func TestParse1(t *testing.T) {
	src := `(1 2 3)`
	st := NewSymbolTable()
	lists, err := ParseString("TestParse1", st, src)
	if err != nil {
		t.Fatalf("Parse error with \"%v\"", err)
	}
	if len(lists) != 1 {
		t.Errorf("Unexpected result")
	}
	lst := lists[0]
	v1, ok := lst.IntAt(0)
	if !ok {
		t.Errorf("Parse error at %v", lst.ElementAt(0).Position())
	}
	if v1 != 1 {
		t.Errorf("Value parse error %d", v1)
	}

	v2, ok := lst.IntAt(1)
	if !ok {
		t.Errorf("Parse error at %v", lst.ElementAt(1).Position())
	}
	if v2 != 2 {
		t.Errorf("Value parse error %d", v2)
	}

	v3, ok := lst.IntAt(2)
	if !ok {
		t.Errorf("Parse error at %v", lst.ElementAt(2).Position())
	}
	if v3 != 3 {
		t.Errorf("Value parse error %d", v3)
	}
}

func TestParse2(t *testing.T) {
	src := `(1.0 2.0 3.0)`
	st := NewSymbolTable()
	lists, err := ParseString("TestParse2", st, src)
	if err != nil {
		t.Fatalf("Parse error with \"%v\"", err)
	}
	if len(lists) != 1 {
		t.Errorf("Unexpected result")
	}
	lst := lists[0]
	v1, ok := lst.FloatAt(0)
	if !ok {
		t.Errorf("Parse error at %v", lst.ElementAt(0).Position())
	}
	if v1 != 1.0 {
		t.Errorf("Value parse error %f", v1)
	}

	v2, ok := lst.FloatAt(1)
	if !ok {
		t.Errorf("Parse error at %v", lst.ElementAt(1).Position())
	}
	if v2 != 2.0 {
		t.Errorf("Value parse error %f", v2)
	}

	v3, ok := lst.FloatAt(2)
	if !ok {
		t.Errorf("Parse error at %v", lst.ElementAt(2).Position())
	}
	if v3 != 3.0 {
		t.Errorf("Value parse error %f", v3)
	}
}

func TestParse3(t *testing.T) {
	src := `("abc" "def" "ghi")`
	st := NewSymbolTable()
	lists, err := ParseString("TestParse3", st, src)
	if err != nil {
		t.Fatalf("Parse error with \"%v\"", err)
	}
	if len(lists) != 1 {
		t.Errorf("Unexpected result")
	}
	lst := lists[0]
	v1, ok := lst.StringAt(0)
	if !ok {
		t.Errorf("Parse error at %v", lst.ElementAt(0).Position())
	}
	if v1 != "abc" {
		t.Errorf("Value parse error %s", v1)
	}

	v2, ok := lst.StringAt(1)
	if !ok {
		t.Errorf("Parse error at %v", lst.ElementAt(1).Position())
	}
	if v2 != "def" {
		t.Errorf("Value parse error %s", v2)
	}

	v3, ok := lst.StringAt(2)
	if !ok {
		t.Errorf("Parse error at %v", lst.ElementAt(2).Position())
	}
	if v3 != "ghi" {
		t.Errorf("Value parse error %s", v3)
	}
}

func TestParse4(t *testing.T) {
	src := `(abc def ghi)`
	st := NewSymbolTable()
	lists, err := ParseString("TestParse3", st, src)
	if err != nil {
		t.Fatalf("Parse error with \"%v\"", err)
	}
	if len(lists) != 1 {
		t.Errorf("Unexpected result")
	}
	lst := lists[0]
	if len(lst.elements) != 3 {
		t.Fatalf("Parse error, unexpected list content")
	}

	strs := []string{"abc", "def", "ghi"}
	for i := 0; i < 3; i++ {
		v, ok := lst.SymbolAt(i)
		if !ok {
			t.Errorf("Parse error at %v", lst.ElementAt(i).Position())
		}
		vs, err := st.GetSymbolName(v)
		if err != nil {
			t.Errorf("Symbol lookup error %d", v)
		}
		if vs != strs[i] {
			t.Errorf("Value parse error %s", vs)
		}
	}
}

func TestParse11(t *testing.T) {
	src := `  
(1 2 3)

(3 4 5)
`
	st := NewSymbolTable()
	vals := [][]int64{{1, 2, 3}, {3, 4, 5}}

	lists, err := ParseString("TestParse11", st, src)
	if err != nil {
		t.Fatalf("Parse error with \"%v\"", err)
	}
	if len(lists) != 2 {
		t.Errorf("Unexpected result")
	}

	for i := 0; i < 2; i++ {
		lst := lists[i]
		for j := 0; j < 3; j++ {
			v, ok := lst.IntAt(j)
			if !ok {
				t.Errorf("Parse error at %v", lst.ElementAt(j).Position())
			}
			if v != vals[i][j] {
				t.Errorf("Value parse error %d", v)
			}
		}
	}
}

func TestParse12(t *testing.T) {
	src := `  ; comment
(1 2 3) ; comment

(3 4 ; comment
5)
`
	st := NewSymbolTable()
	vals := [][]int64{{1, 2, 3}, {3, 4, 5}}

	lists, err := ParseString("TestParse12", st, src)
	if err != nil {
		t.Fatalf("Parse error with \"%v\"", err)
	}
	if len(lists) != 2 {
		t.Errorf("Unexpected result")
	}

	for i := 0; i < 2; i++ {
		lst := lists[i]
		for j := 0; j < 3; j++ {
			v, ok := lst.IntAt(j)
			if !ok {
				t.Errorf("Parse error at %v", lst.ElementAt(j).Position())
			}
			if v != vals[i][j] {
				t.Errorf("Value parse error %d", v)
			}
		}
	}
}

func TestParse13(t *testing.T) {
	src := `(1 2 3`

	st := NewSymbolTable()
	_, err := ParseString("TestParse13", st, src)
	if err == nil {
		t.Error("No parse error")
	}
	pe, ok := err.(*ParseError)
	if !ok {
		if pe.ID != ErrorMissingClosingParenthesis {
			t.Errorf("Unexpected message id: %d", pe.ID)
		}
	}
}
