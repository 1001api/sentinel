package services

import (
	"context"
	"errors"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func initTest(t *testing.T) (*mocks.UserRepo, UserServiceImpl) {
	var mockUserRepo = mocks.NewUserRepo(t)
	var mockIPRepo = mocks.NewIPDBRepo(t)
	var utilService = InitUtilService(&validator.Validate{}, mockIPRepo)
	return mockUserRepo, InitUserService(&utilService, mockUserRepo)
}

func TestFindUserByEmail(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		expectedResult *gen.FindUserByEmailRow
		expectedError  error
	}{
		{
			name:  "Should return user account if email exist",
			email: "correct@gmail.com",
			expectedResult: &gen.FindUserByEmailRow{
				ID:       uuid.MustParse(faker.UUIDHyphenated()),
				Fullname: faker.Name(),
				Email:    "correct@gmail.com",
			},
			expectedError: nil,
		},
		{
			name:           "Should return error if email does not exist",
			email:          "notfound@gmail.com",
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepo, userService := initTest(t)

			var mockReturn gen.FindUserByEmailRow
			if test.expectedResult != nil {
				mockReturn = gen.FindUserByEmailRow{
					ID:       test.expectedResult.ID,
					Fullname: test.expectedResult.Fullname,
					Email:    test.expectedResult.Email,
				}
			}

			// mock repo call
			mockUserRepo.Mock.On("FindUserByEmail", context.Background(), test.email).Return(mockReturn, test.expectedError)

			result, err := userService.FindByEmail(test.email)

			if test.expectedError != nil {
				assert.Nil(t, result)
				assert.NotNil(t, err)
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, result)
				assert.IsType(t, &gen.FindUserByEmailRow{}, result)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}
}

func TestFindUserByEmailWithHash(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		expectedResult *gen.FindUserByEmailWithHashRow
		expectedError  error
	}{
		{
			name:  "Should return user account with password hash if email exist",
			email: "correct@gmail.com",
			expectedResult: &gen.FindUserByEmailWithHashRow{
				ID:             uuid.MustParse(faker.UUIDHyphenated()),
				PasswordHashed: faker.Password(),
				Email:          "correct@gmail.com",
			},
			expectedError: nil,
		},
		{
			name:           "Should return error if email does not exist",
			email:          "notfound@gmail.com",
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepo, userService := initTest(t)

			var mockReturn gen.FindUserByEmailWithHashRow
			if test.expectedResult != nil {
				mockReturn = gen.FindUserByEmailWithHashRow{
					ID:             test.expectedResult.ID,
					PasswordHashed: test.expectedResult.PasswordHashed,
					Email:          test.expectedResult.Email,
				}
			}

			// mock repo call
			mockUserRepo.Mock.On("FindUserByEmailWithHash", context.Background(), test.email).Return(mockReturn, test.expectedError)

			result, err := userService.FindByEmailWithHash(test.email)

			if test.expectedError != nil {
				assert.Nil(t, result)
				assert.NotNil(t, err)
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, result)
				assert.IsType(t, &gen.FindUserByEmailWithHashRow{}, result)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}
}

func TestFindUserByID(t *testing.T) {
	tests := []struct {
		name           string
		id             uuid.UUID
		expectedResult *gen.FindUserByIDRow
		expectedError  error
	}{
		{
			name: "Should return user account if ID exist",
			id:   uuid.MustParse("ef687f52-2c0f-4ba4-ad49-7bfe5df8977b"),
			expectedResult: &gen.FindUserByIDRow{
				ID:       uuid.MustParse("ef687f52-2c0f-4ba4-ad49-7bfe5df8977b"),
				Fullname: faker.Name(),
				Email:    faker.Email(),
			},
			expectedError: nil,
		},
		{
			name:           "Should return error if ID does not exist",
			id:             uuid.MustParse("491d660e-5b7c-4838-9d64-ca70eda51e18"),
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepo, userService := initTest(t)

			var mockReturn gen.FindUserByIDRow
			if test.expectedResult != nil {
				mockReturn = gen.FindUserByIDRow{
					ID:       test.expectedResult.ID,
					Fullname: test.expectedResult.Fullname,
					Email:    test.expectedResult.Email,
				}
			}

			// mock repo call
			mockUserRepo.Mock.On(
				"FindUserByID",
				context.Background(),
				test.id,
			).Return(mockReturn, test.expectedError)

			result, err := userService.FindByID(test.id)

			if test.expectedError != nil {
				assert.Nil(t, result)
				assert.NotNil(t, err)
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, result)
				assert.IsType(t, &gen.FindUserByIDRow{}, result)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}
}

func TestFindByPublicKey(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		expectedResult *gen.FindUserByPublicKeyRow
		expectedError  error
	}{
		{
			name: "Should return user account if key exist",
			key:  faker.Password(),
			expectedResult: &gen.FindUserByPublicKeyRow{
				ID:       uuid.MustParse(faker.UUIDHyphenated()),
				Fullname: faker.Name(),
				Email:    faker.Email(),
			},
			expectedError: nil,
		},
		{
			name:           "Should return error if key does not exist",
			key:            faker.Password(),
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepo, userService := initTest(t)

			var mockReturn gen.FindUserByPublicKeyRow
			if test.expectedResult != nil {
				mockReturn = gen.FindUserByPublicKeyRow{
					ID:       test.expectedResult.ID,
					Fullname: test.expectedResult.Fullname,
					Email:    test.expectedResult.Email,
				}
			}

			// mock repo call
			mockUserRepo.Mock.On(
				"FindUserByPublicKey",
				context.Background(),
				test.key,
			).Return(mockReturn, test.expectedError)

			result, err := userService.FindByPublicKey(test.key)

			if test.expectedError != nil {
				assert.Nil(t, result)
				assert.NotNil(t, err)
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, result)
				assert.IsType(t, &gen.FindUserByPublicKeyRow{}, result)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}
}

func TestGetPublicKey(t *testing.T) {
	tests := []struct {
		name           string
		userID         uuid.UUID
		expectedResult string
		expectedError  error
	}{
		{
			name:           "Should return user public key if user exist",
			userID:         uuid.MustParse(faker.UUIDHyphenated()),
			expectedResult: faker.Password(),
			expectedError:  nil,
		},
		{
			name:           "Should return error if user does not exist",
			userID:         uuid.MustParse(faker.UUIDHyphenated()),
			expectedResult: "",
			expectedError:  errors.New("user not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepo, userService := initTest(t)

			// mock repo call
			mockUserRepo.Mock.On(
				"FindUserPublicKey",
				context.Background(),
				test.userID,
			).Return(test.expectedResult, test.expectedError)

			result, err := userService.GetPublicKey(test.userID)

			if test.expectedError != nil {
				assert.Empty(t, result)
				assert.NotNil(t, err)
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, result)
				assert.IsType(t, "", result)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}
}

func TestCheckAdminExist(t *testing.T) {
	tests := []struct {
		name           string
		expectedResult bool
		expectedError  error
	}{
		{
			name:           "Should return true if admin user exist",
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:           "Should return false if admin user does not exist",
			expectedResult: false,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepo, userService := initTest(t)

			// mock repo call
			mockUserRepo.Mock.On(
				"CheckAdminExist",
				context.Background(),
			).Return(test.expectedResult, test.expectedError)

			result, err := userService.CheckAdminExist()

			if test.expectedError != nil {
				assert.False(t, result)
				assert.NotNil(t, err)
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.True(t, result)
				assert.IsType(t, true, result)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}

}

func TestCreateUser(t *testing.T) {
	expectedFullname := faker.Name()
	expectedEmail := faker.Email()

	tests := []struct {
		name           string
		payload        *gen.CreateUserParams
		expectedResult *gen.CreateUserRow
		expectedError  error
	}{
		{
			name: "Should return true if admin user exist",
			payload: &gen.CreateUserParams{
				Fullname: expectedFullname,
				Email:    expectedEmail,
			},
			expectedResult: &gen.CreateUserRow{
				Fullname: expectedFullname,
				Email:    expectedEmail,
			},
			expectedError: nil,
		},
		{
			name: "Should return false if admin user does not exist",
			payload: &gen.CreateUserParams{
				Fullname: expectedFullname,
				Email:    expectedEmail,
			},
			expectedResult: nil,
			expectedError:  errors.New("something bad happen"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserRepo, userService := initTest(t)

			var mockReturn gen.CreateUserRow
			if test.expectedResult != nil {
				mockReturn = gen.CreateUserRow{
					Fullname: test.expectedResult.Fullname,
					Email:    test.expectedResult.Email,
				}
			}

			// mock repo call
			mockUserRepo.Mock.On(
				"CreateUser",
				context.Background(),
				test.payload,
			).Return(mockReturn, test.expectedError)

			result, err := userService.CreateUser(test.payload)

			if test.expectedError != nil {
				assert.Nil(t, result)
				assert.NotNil(t, err)
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.IsType(t, &gen.CreateUserRow{}, result)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}

}

