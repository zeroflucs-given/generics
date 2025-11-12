package generics_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeroflucs-given/generics"
)

func TestMin(t *testing.T) {
	// Arrange
	input := []int{10, 14, 28, 5}
	target := 5

	// Act
	result := generics.Min(input)

	// Assert
	require.Equal(t, result, target, "Should have the correct value")
}

func TestMinEmpty(t *testing.T) {
	// Arrange
	var input []int
	target := 0

	// Act
	result := generics.Min(input)

	// Assert
	require.Equal(t, result, target, "Should have the correct value")
}

func TestMax(t *testing.T) {
	// Arrange
	input := []int{10, 14, 28, 5}
	target := 28

	// Act
	result := generics.Max(input)

	// Assert
	require.Equal(t, result, target, "Should have the correct value")
}

func TestMaxEmpty(t *testing.T) {
	// Arrange
	var input []int
	target := 0

	// Act
	result := generics.Max(input)

	// Assert
	require.Equal(t, result, target, "Should have the correct value")
}

func TestSum(t *testing.T) {
	// Arrange
	input := []int{10, 14, 28}
	target := 52

	// Act
	result := generics.Sum(input)

	// Assert
	require.Equal(t, result, target, "Should have the correct value")
}
