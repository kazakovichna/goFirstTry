package service

import "testing"

func TestAuthService_PlusSomething(t *testing.T) {
	inValue := 1
	got := PlusSomething(inValue)
	if got != 2 {
		t.Errorf("s.PlusSomething(1) = %d; want 2", got)
	}
}
