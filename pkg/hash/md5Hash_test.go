package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMd5Hash(t *testing.T) {
	// Create some test cases with inputs and expected outputs
	testCases := []struct {
		name   string
		text   string
		result string
	}{
		{
			name:   "empty string",
			text:   "",
			result: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:   "hello world",
			text:   "hello world",
			result: "5eb63bbbe01eeed093cb22bb8f5acdc3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the GetMd5Hash function with the test input
			res := GetMd5Hash(tc.text)

			// Compare the result with the expected output
			assert.Equal(t, tc.result, res)
		})
	}
}
