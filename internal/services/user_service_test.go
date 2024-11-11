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
	mockUserRepo, userService := initTest(t)

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
	mockUserRepo, userService := initTest(t)

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
