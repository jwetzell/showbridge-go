package framer_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/framer"
)

func TestNilGetFramer(t *testing.T) {
	nilFramer := framer.GetFramer("asldfiudchuehrkbjbkjrbb")

	if nilFramer != nil {
		t.Errorf("Expected nil framer, got %v", nilFramer)
	}
}
