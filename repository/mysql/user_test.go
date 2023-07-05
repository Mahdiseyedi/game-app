package mysql

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestIsPhoneNumberUnique(t *testing.T) {
	// Create a mock database object using sqlmock
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Create a DB object with the mock database object
	db := &DB{db: mockDB}

	// Define some test cases with inputs and expected outputs
	testCases := []struct {
		input    string
		expected bool
	}{
		{"09123456789", false}, // phone number exists in the database
		{"09123456780", true},  // phone number does not exist in the database
	}

	// Loop over the test cases and check the output of IsPhoneNumberUnique method
	for _, tc := range testCases {
		// Expect a query to be executed with the mock database object and return a fake row or an error
		if tc.expected {
			mock.ExpectQuery(`select \* from users where phone_number=\?`).WithArgs(tc.input).WillReturnError(sql.ErrNoRows)
		} else {
			mock.ExpectQuery(`select \* from users where phone_number=\?`).WithArgs(tc.input).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone_number", "created_at"}).AddRow(1, "Alice", tc.input, []uint8{0}))
		}

		// Call the IsPhoneNumberUnique method and get the output
		output, _ := db.IsPhoneNumberUnique(tc.input)

		// Check if the output matches the expected value
		if output != tc.expected {
			t.Errorf("IsPhoneNumberUnique(%s) = %v; want %v", tc.input, output, tc.expected)
		}

		// Check if there were any errors during the execution
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
