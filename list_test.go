package listparser

import "testing"

func TestList1(t *testing.T) {
	src := `(event test1 goto state1)`
	st := NewSymbolTable()
	symEvent := st.GetSymbolID("event")
	symGoto := st.GetSymbolID("goto")
	lists, err := ParseString("test1", st, src)
	if err != nil {
		t.Fatalf("Parse error with \"%v\"", err)
	}
	if len(lists) != 1 || lists[0].Len() != 4 {
		t.Error("Parse error")
	}
	if !lists[0].Matches(symEvent, Symbol, symGoto, Symbol) {
		t.Error("Not matched.")
	}
}

func TestList2(t *testing.T) {
	src := `(event test1 goto state1)`
	st := NewSymbolTable()
	symEvent := st.GetSymbolID("event")
	symGoto := st.GetSymbolID("goto")
	lists, err := ParseString("test1", st, src)
	if err != nil {
		t.Fatalf("Parse error with \"%v\"", err)
	}
	if len(lists) != 1 || lists[0].Len() != 4 {
		t.Error("Parse error")
	}
	if lists[0].Matches(symEvent, Int, symGoto, Symbol) {
		t.Error("Matched.")
	}
}

func TestList3(t *testing.T) {
	src := `(event test1 goto state1)`
	st := NewSymbolTable()
	symEvent := st.GetSymbolID("event")
	st.GetSymbolID("goto")
	lists, err := ParseString("test1", st, src)
	if err != nil {
		t.Fatalf("Parse error with \"%v\"", err)
	}
	if len(lists) != 1 || lists[0].Len() != 4 {
		t.Error("Parse error")
	}
	if !lists[0].StartWith(symEvent, symbol) {
		t.Error("Not matched.")
	}
}
