package listparser

import "testing"

func TestList1(t *testing.T) {
	src := `(event test1 goto state1)`
	st := NewSymbolTable()
	symEvent := st.GetSymbolID("event")
	symGoto := st.GetSymbolID("goto")
	lists, err := ParseString("test1", st, src, true, false)
	if err != nil {
		t.Fatalf("Parse error with \"%v\"", err)
	}
	if len(lists) != 1 || lists[0].Len() != 4 {
		t.Error("Parse error")
	}
	lst := lists[0]
	if !(IsSymbolID(lst.ElementAt(0), symEvent) &&
		IsSymbol(lst.ElementAt(1)) &&
		IsSymbolID(lst.ElementAt(2), symGoto) &&
		IsSymbol(lst.ElementAt(3))) {
		t.Error("Not matched.")
	}
}

func TestList2(t *testing.T) {
	src := `(event test1 goto state1)`
	st := NewSymbolTable()
	symEvent := st.GetSymbolID("event")
	symGoto := st.GetSymbolID("goto")
	lists, err := ParseString("test1", st, src, true, false)
	if err != nil {
		t.Fatalf("Parse error with \"%v\"", err)
	}
	if len(lists) != 1 || lists[0].Len() != 4 {
		t.Error("Parse error")
	}
	lst := lists[0]
	if IsSymbolID(lst.ElementAt(0), symEvent) &&
		IsInt(lst.ElementAt(1)) &&
		IsSymbolID(lst.ElementAt(2), symGoto) &&
		IsSymbol(lst.ElementAt(3)) {
		t.Error("Matched.")
	}
}

func TestList3(t *testing.T) {
	src := `(event test1 goto state1)`
	st := NewSymbolTable()
	symEvent := st.GetSymbolID("event")
	st.GetSymbolID("goto")
	lists, err := ParseString("test1", st, src, true, false)
	if err != nil {
		t.Fatalf("Parse error with \"%v\"", err)
	}
	if len(lists) != 1 || lists[0].Len() != 4 {
		t.Error("Parse error")
	}
	lst := lists[0]
	if !(IsSymbolID(lst.ElementAt(0), symEvent) &&
		IsSymbol(lst.ElementAt(1))) {
		t.Error("Not matched.")
	}
}
