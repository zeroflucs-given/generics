package generics_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeroflucs-given/generics"
)

func TestKeyValuesToMap(t *testing.T) {
	// Arrange
	input := []generics.KeyValuePair[string, string]{
		{
			Key:   "hello",
			Value: "world",
		},
		{
			Key:   "fizz",
			Value: "buzz",
		},
	}

	// Act
	result := generics.KeyValuesToMap(input)

	// Assert
	require.NotNil(t, result, "Should have a result")
	require.Len(t, result, 2, "Should have right number of keys")
	require.Equal(t, map[string]string{
		"hello": "world",
		"fizz":  "buzz",
	}, result, "Should have expected output")
}

func TestKeys(t *testing.T) {
	// Arrange
	input := map[string]string{
		"hello": "world",
		"fizz":  "buzz",
	}

	// Act
	result := generics.Keys(input)

	// Assert
	require.NotNil(t, result, "Should have a result")
	require.ElementsMatch(t, result, []string{"hello", "fizz"})
}

func TestValues(t *testing.T) {
	// Arrange
	input := map[string]string{
		"hello": "world",
		"fizz":  "buzz",
	}

	// Act
	result := generics.Values(input)

	// Assert
	require.NotNil(t, result, "Should have a result")
	require.ElementsMatch(t, result, []string{"world", "buzz"})
}

func TestToKeyValues(t *testing.T) {
	// Arrange
	input := map[string]string{
		"hello": "world",
		"fizz":  "buzz",
	}

	// Act
	result := generics.ToKeyValues(input)

	// Assert
	require.NotNil(t, result, "Should have a result")
	require.ElementsMatch(t, result, []generics.KeyValuePair[string, string]{
		{
			Key:   "hello",
			Value: "world",
		},
		{
			Key:   "fizz",
			Value: "buzz",
		},
	})
}
