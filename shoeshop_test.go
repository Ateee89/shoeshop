package main

import "testing"

func TestGenerateID(t *testing.T) {
	testCases := []struct {
		name     string
		input    Shoe
		expected string
	}{
		{"first", Shoe{"id", "model", "brand", 32}, "bra-1"},
		{"second", Shoe{"id", "model", "brand", 32}, "bra-2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nextTc := tc
			actual := GenerateID(nextTc.input)
			if actual != nextTc.expected {
				t.Errorf("test: %q expected: %q actual: %q", nextTc.name, nextTc.expected, actual)
			}
		})
	}
}
