package question

import (
	"testing"
)

func TestPossibleAnswerChoiceIsValid(t *testing.T) {
	// Define some test cases with inputs and expected outputs
	testCases := []struct {
		input    PossibleAnswerChoice
		expected bool
	}{
		{PossibleAnswerA, true}, // valid choice
		{PossibleAnswerB, true}, // valid choice
		{PossibleAnswerC, true}, // valid choice
		{PossibleAnswerD, true}, // valid choice
		{0, false},              // invalid choice
		{5, false},              // invalid choice
	}

	// Loop over the test cases and check the output of IsValid method
	for _, tc := range testCases {
		output := tc.input.IsValid()
		if output != tc.expected {
			t.Errorf("IsValid(%d) = %v; want %v", tc.input, output, tc.expected)
		}
	}
}
