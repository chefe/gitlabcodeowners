// Package testhelper contains utilities for testing.
package testhelper

import (
	"testing"

	"github.com/go-test/deep"
)

// DeepEqual compares the `actual` value against the `expected` value
// using a deep comparison. If there are differences it reports the
// differences and then fails the test.
func DeepEqual(t *testing.T, actual, expected any) {
	t.Helper()

	deep.CompareUnexportedFields = true

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Errorf("values are not equal, excepted=%v, actual=%v, diff=%v", expected, actual, diff)
	}
}
