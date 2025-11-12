package filtering_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeroflucs-given/generics/filtering"
)

// TestWrapContect checks we can wrap a filter with a context awareness
func TestWrapContext(t *testing.T) {
	// Arrange
	ctx := context.Background()
	base := filtering.True[string]

	// Act
	wrapped := filtering.WrapWithContext(base)
	result, err := wrapped(ctx, 1, "hello")

	// Assert
	require.NoError(t, err, "Should not error")
	require.True(t, result, "Should always be true")
}
