package app

import "testing"

func TestHello(t *testing.T) {
	got := Hello()
	if got != nil {
		t.Errorf("Hello() = %d; want nil", got)
	}
}
