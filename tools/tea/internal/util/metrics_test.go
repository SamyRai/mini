package util

import "testing"

func TestUniqStrings(t *testing.T) {
    in := []string{"a", "b", "a", "c", "b"}
    got := UniqStrings(in)
    if len(got) != 3 {
        t.Fatalf("expected 3 unique, got %d: %v", len(got), got)
    }
}
