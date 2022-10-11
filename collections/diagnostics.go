package collections

import "io"

// Diagnosable is an interface that describes the
type Diagnosable interface {
	// Check consistency
	CheckConsistency()

	// Dump the content of a tree (Diagnostics)
	Dump(file io.Writer)
}
