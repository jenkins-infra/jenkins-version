package version_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/garethjevans/jenkins-version/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestByteCountBinary(t *testing.T) {
	type test struct {
		numberOfBytes int64
		expected      string
	}

	tests := []test{
		{numberOfBytes: 0, expected: "0 B"},
		{numberOfBytes: 1, expected: "1 B"},
		{numberOfBytes: 2, expected: "2 B"},
		{numberOfBytes: 1 * 1024, expected: "1.0 KiB"},
		{numberOfBytes: 2 * 1024, expected: "2.0 KiB"},
		{numberOfBytes: 1 * 1024 * 1024, expected: "1.0 MiB"},
		{numberOfBytes: 2 * 1024 * 1024, expected: "2.0 MiB"},
		{numberOfBytes: 1 * 1024 * 1024 * 1024, expected: "1.0 GiB"},
		{numberOfBytes: 2 * 1024 * 1024 * 1024, expected: "2.0 GiB"},
		{numberOfBytes: 27, expected: "27 B"},
		{numberOfBytes: 999, expected: "999 B"},
		{numberOfBytes: 1000, expected: "1000 B"},
		{numberOfBytes: 1023, expected: "1023 B"},
		{numberOfBytes: 1728, expected: "1.7 KiB"},
		{numberOfBytes: 1855425871872, expected: "1.7 TiB"},
		{numberOfBytes: math.MaxInt64, expected: "8.0 EiB"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("TestByteCountBinary-%d", tc.numberOfBytes), func(t *testing.T) {
			o := version.ByteCountBinary(tc.numberOfBytes)
			assert.Equal(t, tc.expected, o)
		})
	}
}
