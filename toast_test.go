package toast

import "testing"

func TestNotification_Push(t *testing.T) {
	n := Notification{
		Title:   "Hello World",
		Message: "It's wild out there, don't forget to take a sword!",
	}
	err := n.Push()
	if err != nil {
		t.Error(err)
	}
}
