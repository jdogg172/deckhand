package unit

import "testing"

func TestSmoke(t *testing.T) {
	if 2+2 != 4 {
		t.Fatal("math broke, pack it up")
	}
}
