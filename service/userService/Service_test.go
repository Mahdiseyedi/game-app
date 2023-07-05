package userService

import (
	"errors"
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

// TestNameServiceValidator tests the NameServiceValidator method of the Service struct
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
		})
	}
}

// TestPhoneNumberServiceValidator tests the PhoneNumberServiceValidator method of the Service struct
//func TestPhoneNumberServiceValidator(t *testing.T) {
//	// Create a mock repository
//	mockRepo := new(MockRepository)
//
//	// Create a service with the mock repository
//	service := New(mockRepo)
//
//	// Create some test cases with inputs and expected outputs
//	testCases := []struct {
//		name    string
//		request RegisterRequest
//		result  bool
//		err     error
//		output  string // expected output from phoneNumber.IsValid function
//	}{
//		{
//			name: "valid phone number",
//			request: RegisterRequest{
//				Name:        "Alice",
//				PhoneNumber: "09123456789", // changed to match phoneNumber.IsValid format
//			},
//			result: true,
//			err:    nil,
//			output: "phone Number is Valid\n", // added expected output from phoneNumber.IsValid function
//		},
//		{
//			name: "invalid phone number",
//			request: RegisterRequest{
//				Name:        "Bob",
//				PhoneNumber: "1234567890",
//			},
//			result: false,
//			err:    errors.New("...Validator: this number is not valid..."),
//			output: "", // no output from phoneNumber.IsValid function
//		},
//		{
//			name: "duplicate phone number",
//			request: RegisterRequest{
//				Name:        "Charlie",
//				PhoneNumber: "09123456789",
//			},
//			result: false,
//			err:    errors.New("...Validator: phone number is not unique..."),
//			output: "phone Number is Valid\n", // added expected output from phoneNumber.IsValid function
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			if tc.err == nil {
//				mockRepo.On("IsPhoneNumberUnique", tc.request.PhoneNumber).Return(true, nil)
//			} else if tc.err.Error() == "...Validator: phone number is not unique..." {
//				mockRepo.On("IsPhoneNumberUnique", tc.request.PhoneNumber).Return(false, nil)
//			}
//
//			res, err := service.PhoneNumberServiceValidator(tc.request)
//
//			assert.Equal(t, tc.result, res)
//			assert.Equal(t, tc.err, err)
//
//			mockRepo.AssertExpectations(t)
//
//			// Create a buffer to capture the output from phoneNumber.IsValid function
//			var buf bytes.Buffer
//
//			// Call the Fprintln function with the buffer and the request phone number
//			fmt.Fprintln(&buf, tc.request.PhoneNumber)
//
//			// Get the output from the buffer
//			output := buf.String()
//
//			assert.Equal(t, tc.output, output)
//		})
//	}
//}
