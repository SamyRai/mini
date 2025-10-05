package graph

import "testing"

func TestFindCycles(t *testing.T) {
	adj := map[string][]string{
		"a": []string{"b"},
		"b": []string{"c"},
		"c": []string{"a"},
	}
	cycles := FindCycles(adj)
	if len(cycles) == 0 {
		t.Fatalf("expected at least one cycle")
	}
}
