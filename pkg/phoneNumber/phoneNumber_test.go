package phoneNumber

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	// Define some test cases with inputs and expected outputs
	testCases := []struct {
		input    string
		expected bool
	}{
		{"09123456789", true},    // valid phone number
		{"+989123456789", false}, // invalid phone number with +98 prefix
		{"0912345678", false},    // invalid phone number with wrong length
		{"01923456789", false},   // invalid phone number with wrong prefix
		{"0912345678a", false},   // invalid phone number with non-digit character
	}

	// Loop over the test cases and check the output of IsValid function
	for _, tc := range testCases {
		output := IsValid(tc.input)
		if output != tc.expected {
			t.Errorf("IsValid(%s) = %v; want %v", tc.input, output, tc.expected)
		}
	}
}
