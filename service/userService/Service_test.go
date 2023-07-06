package userService

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"

	"game-app/entity/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

// IsPhoneNumberUnique is a mock method that returns the arguments passed to it
func (m *MockRepository) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	args := m.Called(phoneNumber)
	return args.Bool(0), args.Error(1)
}

// RegisterUser is a mock method that returns the arguments passed to it
func (m *MockRepository) RegisterUser(u user.User) (user.User, error) {
	args := m.Called(u)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockRepository) GetUserByPhoneNumber(phoneNumber string) (user.User, error) {
	args := m.Called(phoneNumber)
	return args.Get(0).(user.User), args.Error(1)
}

func TestNameServiceValidator(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)

	// Create a service with the mock repository
	service := New(mockRepo)

	// Create some test cases with inputs and expected outputs
	testCases := []struct {
		name    string
		request RegisterRequest
		result  bool
		err     error
	}{
		{
			name: "valid name",
			request: RegisterRequest{
				Name:        "Alice",
				PhoneNumber: "+1234567890",
			},
			result: true,
			err:    nil,
		},
		{
			name: "short name",
			request: RegisterRequest{
				Name:        "Ed",
				PhoneNumber: "+1987654321",
			},
			result: false,
			err:    errors.New("...Validator: name lenght should grater than 3"),
		},
		{
			name: "reserved name",
			request: RegisterRequest{
				Name:        "userName",
				PhoneNumber: "+1234567891",
			},
			result: false,
			err:    errors.New("...Validator: username cant be \"userName\""),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := service.NameServiceValidator(tc.request)

			assert.Equal(t, tc.result, res)
			assert.Equal(t, tc.err, err)

			// No need to call mockRepo.AssertExpectations(t) since no mock methods are used in this test
		})
	}
}

// TestPhoneNumberServiceValidator tests the PhoneNumberServiceValidator method of the Service struct
func TestPhoneNumberServiceValidator(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)

	// Create a service with the mock repository
	service := New(mockRepo)

	// Create some test cases with inputs and expected outputs
	testCases := []struct {
		name    string
		request RegisterRequest
		result  bool
		err     error
	}{
		{
			name: "valid phone number",
			request: RegisterRequest{
				Name:        "Alice",
				PhoneNumber: "09129156789", // removed the + sign
			},
			result: true,
			err:    nil,
		},
		{
			name: "invalid phone number",
			request: RegisterRequest{
				Name:        "Bob",
				PhoneNumber: "1234567890",
			},
			result: false,
			err:    errors.New("...Validator: this number is not valid..."),
		},
		{
			name: "duplicate phone number",
			request: RegisterRequest{
				Name:        "Charlie",
				PhoneNumber: "09123456789", // removed the + sign
			},
			result: false,
			err:    errors.New("...Validator: phone number is not unique..."),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				mockRepo.On("IsPhoneNumberUnique", tc.request.PhoneNumber).Return(true, nil)
			} else if tc.err.Error() == "...Validator: phone number is not unique..." {
				mockRepo.On("IsPhoneNumberUnique", tc.request.PhoneNumber).Return(false, nil)
			}

			res, err := service.PhoneNumberServiceValidator(tc.request)

			assert.Equal(t, tc.result, res)
			assert.Equal(t, tc.err, err)

			mockRepo.AssertExpectations(t)

		})
	}
}

// TestPasswordServiceValidator tests the PasswordServiceValidator method of the Service struct
func TestPasswordServiceValidator(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)

	// Create a service with the mock repository
	service := New(mockRepo)

	// Create some test cases with inputs and expected outputs
	testCases := []struct {
		name    string
		request RegisterRequest
		result  bool
		err     error
	}{
		{
			name: "valid password",
			request: RegisterRequest{
				Name:        "Alice",
				PhoneNumber: "09123456789",
				Password:    "secret123", // valid password with at least 8 characters and not equal to "password"
			},
			result: true,
			err:    nil,
		},
		{
			name: "short password",
			request: RegisterRequest{
				Name:        "Bob",
				PhoneNumber: "09123456789",
				Password:    "pass", // invalid password with less than 8 characters
			},
			result: false,
			err:    errors.New("...Validator: Password len most grater than 8..."),
		},
		{
			name: "simple password",
			request: RegisterRequest{
				Name:        "Charlie",
				PhoneNumber: "09123456789",
				Password:    "password", // invalid password equal to "password"
			},
			result: false,
			err:    errors.New("...Validator: so simple..."),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := service.PasswordServiceValidator(tc.request)

			assert.Equal(t, tc.result, res)
			assert.Equal(t, tc.err, err)

		})
	}
}

func TestRegister(t *testing.T) {
	// Create a mock repository
	mockRepo := new(MockRepository)

	// Create a service with the mock repository
	service := New(mockRepo)

	// Create some test cases with inputs and expected outputs
	testCases := []struct {
		name    string
		request RegisterRequest
		result  RegisterResponse
		err     error
	}{
		{
			name: "valid registration",
			request: RegisterRequest{
				Name:        "Alice",
				PhoneNumber: "09786456789",
				Password:    "secret123",
			},
			result: RegisterResponse{
				user: user.User{
					ID:          1, // mock ID returned by the repository
					Name:        "Alice",
					PhoneNumber: "09198996789",
					Password:    "5ebe2294ecd0e0f08eab7690d2a6ee69", // hashed password using md5
				},
			},
			err: nil,
		},

		{
			name: "duplicate phone number",
			request: RegisterRequest{
				Name:        "Charlie",
				PhoneNumber: "09123456789", // same as Alice's phone number
				Password:    "secret789",
			},
			result: RegisterResponse{},
			err:    errors.New("...Validator: phone number is not unique..."),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				mockRepo.On("IsPhoneNumberUnique", tc.request.PhoneNumber).Return(true, nil)
				mockRepo.On("RegisterUser", mock.AnythingOfType("user.User")).Return(tc.result.user, nil)
			} else if tc.err.Error() == "...Validator: this number is not valid..." {
				mockRepo.On("IsPhoneNumberUnique", tc.request.PhoneNumber).Return(false, errors.New("phone number is invalid"))
			} else if tc.err.Error() == "...Validator: phone number is not unique..." {
				mockRepo.On("IsPhoneNumberUnique", tc.request.PhoneNumber).Return(false, nil)
			}

			res, err := service.Register(tc.request)

			assert.Equal(t, tc.result, res)
			assert.Equal(t, tc.err, err)

			mockRepo.AssertExpectations(t)

		})
	}
}

//-----------------------------------------------------------------------------------------------------

// MockDB is a mock implementation of the DB interface
type MockDB struct {
	mock.Mock
}

// Close is a mock method that returns the arguments passed to it
func (m *MockDB) Close() error {
	mockArgs := m.Called()
	return mockArgs.Error(0)
}

func (m *MockDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	args := m.Called(phoneNumber)
	return args.Bool(0), args.Error(1)
}

// RegisterUser is a mock method that returns the arguments passed to it
func (m *MockDB) RegisterUser(u user.User) (user.User, error) {
	args := m.Called(u)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockDB) GetUserByPhoneNumber(phoneNumber string) (user.User, error) {
	args := m.Called(phoneNumber)
	return args.Get(0).(user.User), args.Error(1)
}

func TestLogin(t *testing.T) {
	// Create a mock database
	mockDB := new(MockDB)

	// Create a repository with the mock database
	//repo := New(mockDB)

	// Create a service with the repository
	service := New(mockDB)

	// Create some test cases with inputs and expected outputs
	testCases := []struct {
		name    string
		request LoginRequest
		result  LoginResponse
		err     error
	}{
		{
			name: "valid login",
			request: LoginRequest{
				PhoneNumber: "09123456789", // same as Alice's phone number
				Password:    "secret123",   // same as Alice's password
			},
			result: LoginResponse{
				user: user.User{
					ID:          1,                                  // same as Alice's ID
					Name:        "Alice",                            // same as Alice's name
					PhoneNumber: "09123456789",                      // same as Alice's phone number
					Password:    "5ebe2294ecd0e0f08eab7690d2a6ee69", // hashed password using md5
				},
			},
			err: nil,
		},

		{
			name: "wrong password",
			request: LoginRequest{
				PhoneNumber: "09123456789", // same as Alice's phone number
				Password:    "wrong123",    // different from Alice's password
			},
			result: LoginResponse{},
			err:    errors.New("...Service: Login failed!..."),
		},

		{
			name: "non-existent phone number",
			request: LoginRequest{
				PhoneNumber: "09111111111", // different from any registered user's phone number
				Password:    "secret123",   // any password
			},
			result: LoginResponse{},
			err:    errors.New("...Repository: Get user by phone number repository Error sql.ErrNoRows"), // assuming this is the error returned by the repository when no user is found
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var rows *sqlmock.Rows

			if tc.err == nil {
				rows = sqlmock.NewRows([]string{"id", "name", "phone_number", "created_at", "password"}).AddRow(tc.result.user.ID, tc.result.user.Name, tc.result.user.PhoneNumber, []uint8{0}, tc.result.user.Password)
			} else if tc.err.Error() == "...Service: Login failed!..." {
				rows = sqlmock.NewRows([]string{"id", "name", "phone_number", "created_at", "password"}).AddRow(1, "Alice", "09123456789", []uint8{0}, "5ebe2294ecd0e0f08eab7690d2a6ee69")
			} else if tc.err.Error() == "...Repository: Get user by phone number repository Error sql.ErrNoRows" {
				rows = sqlmock.NewRows([]string{"id", "name", "phone_number", "created_at", "password"}).RowError(0, sql.ErrNoRows)
			}

			mockDB.On("QueryRow", mock.AnythingOfType("string"), tc.request.PhoneNumber).Return(rows, tc.err)

			res, err := service.Login(tc.request)

			assert.Equal(t, tc.result, res)
			assert.Equal(t, tc.err, err)

			mockDB.AssertExpectations(t)

		})
	}
}
