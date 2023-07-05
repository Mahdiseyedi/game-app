package mysql

import (
	"database/sql"
	"errors"
	"game-app/entity/user"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
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

func TestRegisterUser(t *testing.T) {
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
		input    user.User
		expected user.User
		err      error
	}{
		{
			input: user.User{
				Name:        "Bob",
				PhoneNumber: "09123456780",
				Password:    "hashed_password",
			},
			expected: user.User{
				ID:          2,
				Name:        "Bob",
				PhoneNumber: "09123456780",
				Password:    "hashed_password",
			},
			err: nil,
		},
		{
			input: user.User{
				Name:        "Charlie",
				PhoneNumber: "09123456789", // duplicate phone number
				Password:    "hashed_password",
			},
			expected: user.User{},
			err:      errors.New("cant inseret into DB, ...some thing went wrong, ..."),
		},
	}

	// Loop over the test cases and check the output of RegisterUser method
	for _, tc := range testCases {
		// Expect an exec to be executed with the mock database object and return a fake result or an error
		if tc.err == nil {
			mock.ExpectExec(`insert into users\(name, phone_number,password\) values\(\?,\?,\?\)`).WithArgs(tc.input.Name, tc.input.PhoneNumber, tc.input.Password).WillReturnResult(sqlmock.NewResult(2, 1))
		} else {
			mock.ExpectExec(`insert into users\(name, phone_number,password\) values\(\?,\?,\?\)`).WithArgs(tc.input.Name, tc.input.PhoneNumber, tc.input.Password).WillReturnError(tc.err)
		}

		// Call the RegisterUser method and get the output
		output, err := db.RegisterUser(tc.input)

		// Check if the output matches the expected value
		if !reflect.DeepEqual(output, tc.expected) {
			t.Errorf("RegisterUser(%v) = %v; want %v", tc.input, output, tc.expected)
		}

		// Check if the error matches the expected value
		if !errors.Is(err, tc.err) {
			t.Errorf("RegisterUser(%v) error = %v; want %v", tc.input, err, tc.err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetUserByPhoneNumber(t *testing.T) {
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
		expected user.User
		err      error
	}{
		{
			input: "09123456789",
			expected: user.User{
				ID:          1,
				Name:        "Alice",
				PhoneNumber: "09123456789",
				Password:    "hashed_password",
			},
			err: nil,
		},
		{
			input:    "09123456780",
			expected: user.User{},
			err:      errors.New("...no user find with that phone number, ..."),
		},
	}

	// Loop over the test cases and check the output of GetUserByPhoneNumber method
	for _, tc := range testCases {
		// Expect a query to be executed with the mock database object and return a fake row or an error
		if tc.err == nil {
			mock.ExpectQuery(`select \* from users where phone_number =\?`).WithArgs(tc.input).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone_number", "created_at", "password"}).AddRow(tc.expected.ID, tc.expected.Name, tc.expected.PhoneNumber, []uint8{0}, tc.expected.Password))
		} else {
			mock.ExpectQuery(`select \* from users where phone_number =\?`).WithArgs(tc.input).WillReturnError(tc.err)
		}
		output, err := db.GetUserByPhoneNumber(tc.input)

		// Check if the output matches the expected value
		if !reflect.DeepEqual(output, tc.expected) {
			t.Errorf("GetUserByPhoneNumber(%s) = %v; want %v", tc.input, output, tc.expected)
		}

		// Check if the error matches the expected value
		if !errors.Is(err, tc.err) {
			t.Errorf("GetUserByPhoneNumber(%s) error = %v; want %v", tc.input, err, tc.err)
		}

		// Check if there were any errors during the execution
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
