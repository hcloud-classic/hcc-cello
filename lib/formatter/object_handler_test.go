package formatter

import (
	"testing"
)

func TestReadupdate(t *testing.T) {
	if (readupdate()) == "No" {
		t.Errorf("Wrong, received")
	} else {
		t.Log()
	}
}
