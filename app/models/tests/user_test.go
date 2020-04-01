package models_test

import (
	"testing"

	"github.com/inhumanLightBackend/app/models"
	"github.com/stretchr/testify/assert"
)

func TestUser_Validations(t *testing.T) {
	testCases := []struct {
		name string
		u func() *models.User
		isValid bool
	}{
		{
			name: "Valid user",
			u: func () *models.User  {
				return models.NewTestUser(t)
			},
			isValid: true,
		},
		{
			name: "Empty email",
			u: func () *models.User {
				user := models.NewTestUser(t)
				user.Email = ""
				return user
			},
			isValid: false,
		},
		{
			name: "Email not valid",
			u: func () *models.User {
				user := models.NewTestUser(t)
				user.Email = "kakaoto email@"
				return user
			},
			isValid: false,
		},
		{
			name: "Password is empty",
			u: func () *models.User {
				user := models.NewTestUser(t)
				user.Password = ""
				return user
			},
			isValid: false,
		},
		{
			name: "invalid password",
			u: func () *models.User {
				user := models.NewTestUser(t)
				user.Password = "123"
				return user
			},
			isValid: false,
		},
		{
			name: "with encrypted",
			u: func () *models.User {
				user := models.NewTestUser(t)
				user.Password = ""
				user.EncryptedPassword = "1234567"
				return user
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	user := models.NewTestUserEmptyFields(t)
	
	assert.NoError(t, user.BeforeCreate())
	assert.NotEmpty(t, user.EncryptedPassword)
	assert.NotEmpty(t, user.Role)
	assert.NotEmpty(t, user.IsActive)
	assert.NotEmpty(t, user.Token)
	assert.NotEmpty(t, user.CreatedAt)
}

func TestUser_SetPassword(t *testing.T) {
	user := models.NewTestUser(t)
	oldPassword := user.EncryptedPassword

	assert.NoError(t, user.SetPassword("6543321"))
	assert.NotEmpty(t, user.EncryptedPassword)
	assert.NotEqual(t, user.EncryptedPassword, oldPassword)
}

func TestUser_ChangeActiveStatus(t *testing.T) {
	user := models.NewTestUser(t)
	oldStatus := user.IsActive
	user.ChangeActiveStatus(false)

	assert.NotEqual(t, user.IsActive, oldStatus)
}

func TestUser_ComparePassword(t *testing.T) {
	user := models.NewTestUser(t)
	assert.Equal(t, !user.ComparePassword("123456"), true)
}

