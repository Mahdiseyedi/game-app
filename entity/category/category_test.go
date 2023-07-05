package category

import (
	"reflect"
	"testing"
)

func TestCategoryIsValid(t *testing.T) {
	// Define some test cases with inputs and expected outputs
	testCases := []struct {
		input    Category
		expected bool
	}{
		{FootballCategory, true}, // valid category
		{HistoryCategory, false}, // invalid category
		{"Music", false},         // invalid category
	}

	// Loop over the test cases and check the output of IsValid method
	for _, tc := range testCases {
		output := tc.input.IsValid()
		if output != tc.expected {
			t.Errorf("IsValid(%s) = %v; want %v", tc.input, output, tc.expected)
		}
	}
}

func TestCategoryList(t *testing.T) {
	// Define the expected output
	expected := []Category{FootballCategory, HistoryCategory}

	// Call the CategoryList function and get the output
	output := CategoryList()

	// Compare the output with the expected value using reflect.DeepEqual
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("CategoryList() = %v; want %v", output, expected)
	}
}
