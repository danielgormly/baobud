package core

import "testing"

func TestAnalyze(t *testing.T) {
	got := -1
	if got != 1 {
		t.Errorf("Abs(-1) = %d; want 1", got)
	}
}
