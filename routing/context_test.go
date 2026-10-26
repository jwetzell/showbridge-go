package routing_test

import (
	"io"
	"testing"

	"github.com/jwetzell/showbridge-go/routing"
	"github.com/stretchr/testify/require"
)

func TestGoodNewContextWriter(t *testing.T) {

	testCases := []struct {
		name   string
		writer io.Writer
	}{
		{
			name:   "dummy writer",
			writer: io.Discard,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			contextWriter := routing.NewContextWriter(testCase.writer)

			contextWriter.Lock()
			lockable := contextWriter.TryLock()
			require.False(t, lockable)
			contextWriter.Unlock()

			require.NotNil(t, contextWriter)
			_, err := contextWriter.Write([]byte("test"))
			require.NoError(t, err)
			require.True(t, contextWriter.Done())
		})
	}
}
