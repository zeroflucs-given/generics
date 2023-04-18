package generics_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeroflucs-given/generics"
)

func TestFalse(t *testing.T) {
	// Arrange
	cond := false
	resTrue := 5
	resFalse := 10

	// Act
	result := generics.If(cond, resTrue, resFalse)

	// Assert
	require.Equal(t, result, resFalse, "Should have the correct value")
}

func TestTrue(t *testing.T) {
	// Arrange
	cond := true
	resTrue := 5
	resFalse := 10

	// Act
	result := generics.If(cond, resTrue, resFalse)

	// Assert
	require.Equal(t, result, resTrue, "Should have the correct value")
}
