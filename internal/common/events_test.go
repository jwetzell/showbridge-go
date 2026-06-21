package common_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
)

func TestEventToJson(t *testing.T) {
	e := common.Event{
		Type: "test.event",
		Data: map[string]any{
			"key": "value",
		},
	}
	jsonData, err := e.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert event to JSON: %v", err)
	}
	expectedJson := `{"type":"test.event","data":{"key":"value"}}`
	if string(jsonData) != expectedJson {
		t.Errorf("Expected JSON: %s, got: %s", expectedJson, string(jsonData))
	}
}
