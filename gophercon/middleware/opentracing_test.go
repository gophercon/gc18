package middleware

import "testing"

func TestOperation(t *testing.T) {
	action := operation("github.com/bketelsen/thing/action.Name")
	if action != "Name" {
		t.Errorf("expected %s, got %s", action, "Name")
	}
}
