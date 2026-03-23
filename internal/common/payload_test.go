package common_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/test"
)

func TestGoodGetWrappedPayload(t *testing.T) {
	testCases := []struct {
		name     string
		ctx      context.Context
		payload  any
		expected common.WrappedPayload
	}{
		{
			name:    "basic",
			ctx:     t.Context(),
			payload: "test",
			expected: common.WrappedPayload{
				Payload: "test",
			},
		},
		{
			name:    "with modules in context",
			ctx:     test.GetContextWithModules(t.Context(), map[string]common.Module{}),
			payload: "test",
			expected: common.WrappedPayload{
				Payload: "test",
				Modules: map[string]common.Module{},
			},
		},
		{
			name:    "with sender in context",
			ctx:     test.GetContextWithSender(t.Context(), "sender"),
			payload: "test",
			expected: common.WrappedPayload{
				Payload: "test",
				Sender:  "sender",
			},
		},
		{
			name:    "with source in context",
			ctx:     test.GetContextWithSource(t.Context(), "source"),
			payload: "test",
			expected: common.WrappedPayload{
				Payload: "test",
				Source:  "source",
			},
		},
		{
			name: "with all fields in context",
			ctx: test.GetContextWithSource(
				test.GetContextWithSender(
					test.GetContextWithModules(t.Context(), map[string]common.Module{}),
					"sender",
				),
				"source",
			),
			payload: "test",
			expected: common.WrappedPayload{
				Payload: "test",
				Modules: map[string]common.Module{},
				Sender:  "sender",
				Source:  "source",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wrappedPayload := common.GetWrappedPayload(testCase.ctx, testCase.payload)

			if !reflect.DeepEqual(wrappedPayload, testCase.expected) {
				t.Fatalf("GetWrappedPayload expected got %+v,  expected %+v", wrappedPayload, testCase.expected)
			}
		})
	}

}
