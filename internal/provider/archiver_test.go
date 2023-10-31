// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CheckMatch(t *testing.T) {
	tests := []struct {
		fileName string
		excludes []string
		expected bool
	}{
		{
			fileName: "foo.txt",
			excludes: []string{"foo.txt"},
			expected: true,
		},
		{
			fileName: "foo.txt",
			excludes: []string{"fo?.txt"},
			expected: true,
		},
		{
			fileName: "foo.txt",
			excludes: []string{"f*.txt"},
			expected: true,
		},
		{
			fileName: "foo.txt",
			excludes: []string{"foo.exe", "bar.txt"},
			expected: false,
		},
		{
			fileName: "foo.txt",
			excludes: []string{"foo.exe", "foo.*"},
			expected: true,
		},
	}

	for _, tt := range tests {
		m, err := checkMatch(tt.fileName, tt.excludes)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, tt.expected, m)
	}
}
