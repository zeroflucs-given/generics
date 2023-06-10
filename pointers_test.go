package generics_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeroflucs-given/generics"
)

func TestPointerTo(t *testing.T) {
	// Arrange
	v := "hello"

	// Act
	pv := generics.PointerTo(v)

	// Assert
	require.NotNil(t, pv)
	require.Equal(t, *pv, "hello")
}

func TestPointerOrNil(t *testing.T) {
	// Arrange
	v := ""

	// Act
	pv := generics.PointerOrNil(v)

	// Assert
	require.Nil(t, pv)
}

func TestPointerOrNilWithValue(t *testing.T) {
	// Arrange
	v := "hello"

	// Act
	pv := generics.PointerOrNil(v)

	// Assert
	require.NotNil(t, pv)
	require.Equal(t, *pv, "hello")
}

func TestValueOrDefault(t *testing.T) {
	// Arrange
	x := 4
	pv := &x

	// Act
	v := generics.ValueOrDefault(pv)

	// Assert
	require.Equal(t, 4, v)
}

func TestValueOrDefaultNil(t *testing.T) {
	// Arrange
	var pv *int

	// Act
	v := generics.ValueOrDefault(pv)

	// Assert
	require.Equal(t, 0, v)
}
