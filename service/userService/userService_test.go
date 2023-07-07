package userService

//
//import (
//	"errors"
//	"testing"
//
//	"game-app/entity/user"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//)
//
//// MockRepository is a mock implementation of the Repository interface
//type MockRepository struct {
//	mock.Mock
//}
//
//// IsPhoneNumberUnique is a mock method that returns the arguments passed to it
//func (m *MockRepository) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
//	args := m.Called(phoneNumber)
//	return args.Bool(0), args.Error(1)
//}
//
//// RegisterUser is a mock method that returns the arguments passed to it
//func (m *MockRepository) RegisterUser(u user.User) (user.User, error) {
//	args := m.Called(u)
//	return args.Get(0).(user.User), args.Error(1)
//}
//
//func (m *MockRepository) GetUserByPhoneNumber(phoneNumber string) (user.User, error) {
//	args := m.Called(phoneNumber)
//	return args.Get(0).(user.User), args.Error(1)
//}
//
//func (m *MockRepository) GetUserByID(userID uint) (user.User, error) {
//	args := m.Called(userID)
//	return args.Get(0).(user.User), args.Error(1)
//}
//
//func TestNameServiceValidator(t *testing.T) {
//	// Create a mock repository
//	mockRepo := new(MockRepository)
//
//	// Create a service with the mock repository
//	service := New(mockRepo, JwtSignKey)
//
//	// Create some test cases with inputs and expected outputs
//	testCases := []struct {
//		name    string
//		request RegisterRequest
//		result  bool
//		err     error
//	}{
//		{
//			name: "valid name",
//			request: RegisterRequest{
//				Name:        "Alice",
//				PhoneNumber: "+1234567890",
//			},
//			result: true,
//			err:    nil,
//		},
//		{
//			name: "short name",
//			request: RegisterRequest{
//				Name:        "Ed",
//				PhoneNumber: "+1987654321",
//			},
//			result: false,
//			err:    errors.New("...Validator: name lenght should grater than 3"),
//		},
//		{
//			name: "reserved name",
//			request: RegisterRequest{
//				Name:        "userName",
//				PhoneNumber: "+1234567891",
//			},
//			result: false,
//			err:    errors.New("...Validator: username cant be \"userName\""),
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			res, err := service.NameServiceValidator(tc.request)
//
//			assert.Equal(t, tc.result, res)
//			assert.Equal(t, tc.err, err)
//
//			// No need to call mockRepo.AssertExpectations(t) since no mock methods are used in this test
//		})
//	}
//}
//
//// TestPhoneNumberServiceValidator tests the PhoneNumberServiceValidator method of the Service struct
//func TestPhoneNumberServiceValidator(t *testing.T) {
//	// Create a mock repository
//	mockRepo := new(MockRepository)
//
//	// Create a service with the mock repository
//	service := New(mockRepo, JwtSignKey)
//
//	// Create some test cases with inputs and expected outputs
//	testCases := []struct {
//		name    string
//		request RegisterRequest
//		result  bool
//		err     error
//	}{
//		{
//			name: "valid phone number",
//			request: RegisterRequest{
//				Name:        "Alice",
//				PhoneNumber: "09129156789", // removed the + sign
//			},
//			result: true,
//			err:    nil,
//		},
//		{
//			name: "invalid phone number",
//			request: RegisterRequest{
//				Name:        "Bob",
//				PhoneNumber: "1234567890",
//			},
//			result: false,
//			err:    errors.New("...Validator: this number is not valid..."),
//		},
//		{
//			name: "duplicate phone number",
//			request: RegisterRequest{
//				Name:        "Charlie",
//				PhoneNumber: "09123456789", // removed the + sign
//			},
//			result: false,
//			err:    errors.New("...Validator: phone number is not unique..."),
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
//		})
//	}
//}
//
//// TestPasswordServiceValidator tests the PasswordServiceValidator method of the Service struct
//func TestPasswordServiceValidator(t *testing.T) {
//	// Create a mock repository
//	mockRepo := new(MockRepository)
//
//	// Create a service with the mock repository
//	service := New(mockRepo, JwtSignKey)
//
//	// Create some test cases with inputs and expected outputs
//	testCases := []struct {
//		name    string
//		request RegisterRequest
//		result  bool
//		err     error
//	}{
//		{
//			name: "valid password",
//			request: RegisterRequest{
//				Name:        "Alice",
//				PhoneNumber: "09123456789",
//				Password:    "secret123", // valid password with at least 8 characters and not equal to "password"
//			},
//			result: true,
//			err:    nil,
//		},
//		{
//			name: "short password",
//			request: RegisterRequest{
//				Name:        "Bob",
//				PhoneNumber: "09123456789",
//				Password:    "pass", // invalid password with less than 8 characters
//			},
//			result: false,
//			err:    errors.New("...Validator: Password len most grater than 8..."),
//		},
//		{
//			name: "simple password",
//			request: RegisterRequest{
//				Name:        "Charlie",
//				PhoneNumber: "09123456789",
//				Password:    "password", // invalid password equal to "password"
//			},
//			result: false,
//			err:    errors.New("...Validator: so simple..."),
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			res, err := service.PasswordServiceValidator(tc.request)
//
//			assert.Equal(t, tc.result, res)
//			assert.Equal(t, tc.err, err)
//
//		})
//	}
//}
//
//func TestRegister(t *testing.T) {
//	// Create a mock repository
//	mockRepo := new(MockRepository)
//
//	// Create a service with the mock repository
//	service := New(mockRepo, JwtSignKey)
//
//	// Create some test cases with inputs and expected outputs
//	testCases := []struct {
//		name    string
//		request RegisterRequest
//		result  RegisterResponse
//		err     error
//	}{
//		{
//			name: "valid registration",
//			request: RegisterRequest{
//				Name:        "Alice",
//				PhoneNumber: "09786456789",
//				Password:    "secret123",
//			},
//			result: RegisterResponse{
//				User: user.User{
//					ID:          1, // mock ID returned by the repository
//					Name:        "Alice",
//					PhoneNumber: "09198996789",
//					Password:    "5ebe2294ecd0e0f08eab7690d2a6ee69", // hashed password using md5
//				},
//			},
//			err: nil,
//		},
//
//		{
//			name: "duplicate phone number",
//			request: RegisterRequest{
//				Name:        "Charlie",
//				PhoneNumber: "09123456789", // same as Alice's phone number
//				Password:    "secret789",
//			},
//			result: RegisterResponse{},
//			err:    errors.New("...Validator: phone number is not unique..."),
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			if tc.err == nil {
//				mockRepo.On("IsPhoneNumberUnique", tc.request.PhoneNumber).Return(true, nil)
//				mockRepo.On("RegisterUser", mock.AnythingOfType("user.User")).Return(tc.result.User, nil)
//			} else if tc.err.Error() == "...Validator: this number is not valid..." {
//				mockRepo.On("IsPhoneNumberUnique", tc.request.PhoneNumber).Return(false, errors.New("phone number is invalid"))
//			} else if tc.err.Error() == "...Validator: phone number is not unique..." {
//				mockRepo.On("IsPhoneNumberUnique", tc.request.PhoneNumber).Return(false, nil)
//			}
//
//			res, err := service.Register(tc.request)
//
//			assert.Equal(t, tc.result, res)
//			assert.Equal(t, tc.err, err)
//
//			mockRepo.AssertExpectations(t)
//
//		})
//	}
//}
//
////-----------------------------------------------------------------------------------------------------
////Login worked well
