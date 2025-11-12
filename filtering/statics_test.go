package filtering_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeroflucs-given/generics/filtering"
)

// TestTrueFalse tests the standard true/false constants
func TestTrueFalse(t *testing.T) {
	// Arrange
	tf := filtering.True[string]
	ff := filtering.False[string]

	// Act
	tv := tf(1, "hello")
	fv := ff(1, "hello")

	// Assert
	require.True(t, tv, "Should get a true")
	require.False(t, fv, "Should get a false")
}

// TestTrueFalseWithContext tests the standard true/false constants with contexts
func TestTrueFalseWithContext(t *testing.T) {
	// Arrange
	ctx := context.Background()
	tf := filtering.TrueWithContext[string]
	ff := filtering.FalseWithContext[string]

	// Act
	tv, errTrue := tf(ctx, 1, "hello")
	fv, errFalse := ff(ctx, 1, "hello")

	// Assert
	require.True(t, tv, "Should get a true")
	require.False(t, fv, "Should get a false")
	require.NoError(t, errTrue, "True constant should not error")
	require.NoError(t, errFalse, "False constant should not error")
}
