package generics_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeroflucs-given/generics"
)

func ExampleMust() {
	// Any function that returns a value/error tuple
	fn := func(fail bool) (string, error) {
		if fail {
			return "", errors.New("boom")
		}
		return "Hello", nil
	}

	// Instead of assigning to two variables, can step straight to value
	v := generics.Must(fn(false))
	fmt.Println(v)
}

func TestMustSuccess(t *testing.T) {
	// Arrange
	fn := func() (string, error) {
		return "Hello", nil
	}

	// Act
	v := generics.Must(fn())

	// Assert
	require.Equal(t, "Hello", v)
}

// TestMustFail checks that we get a panic when Must is given an error argument
func TestMustFail(t *testing.T) {
	// Arrange
	fn := func() (string, error) {
		return "", errors.New("failure")
	}
	defer func() {
		rec := recover()
		if rec == nil {
			t.Error("Should get a panic, but did not")
		} else {
			t.Logf("Handled panic: %q", rec)
		}
	}()

	// Act
	v := generics.Must(fn())

	// Assert
	t.Errorf("Should not get here. Somehow got %v", v)
}

func TestValueOrError(t *testing.T) {
	v1, err1 := generics.ValueOrError(64, nil)
	require.Equal(t, 64, v1, "Should get value")
	require.NoError(t, err1, "Should have no error")

	v2, err2 := generics.ValueOrError(64, io.EOF)
	require.Error(t, err2, "Should have the error")
	require.Zero(t, v2, "Should have a zero-value for the type back")
}
