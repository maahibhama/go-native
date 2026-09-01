package ui

import "testing"

type testScheduler struct{ calls int }

func (s *testScheduler) Schedule() { s.calls++ }
func TestState(t *testing.T) {
	sc := &testScheduler{}
	SetScheduler(sc)
	defer SetScheduler(nil)
	s := NewState(1)
	s.Update(func(v int) int { return v + 1 })
	if s.Get() != 2 || sc.calls != 1 {
		t.Fatalf("value=%d calls=%d", s.Get(), sc.calls)
	}
}
