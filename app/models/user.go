package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int        `json:"id"`
	Login             string     `json:"login"`
	Email             string     `json:"email"`
	Password          string     `json:"password,omitempty"`
	EncryptedPassword string     `json:"-"`
	CreatedAt         time.Time `json:"registration_date"`
	Token             string     `json:"api_token"`
	Contacts          string     `json:"contacts"`
	Role              string     `json:"user_role"`
	IsActive          bool       `json:"-"`
}

var (
	Roles = []string{"USER", "ADMIN"}
)

func (user *User) Validate() error {
	return validation.ValidateStruct(
		user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.By(requiredIf(user.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

func (user *User) BeforeCreate() error {
	if len(user.Password) > 0 {
		enc, err := encryptString(user.Password)
		if err != nil {
			return err
		}
		user.EncryptedPassword = enc
	}

	user.CreatedAt = time.Now()
	user.Token = generateToken(user)
	user.Role = Roles[0]
	user.IsActive = true

	return nil
}

func (user *User) SetPassword(password string) error {
	if len(password) > 0 {
		enc, err := encryptString(password)
		if err != nil {
			return err
		}
		user.EncryptedPassword = enc
		return nil
	}

	return errors.New("Empty param: 'password'")
}

func (user *User) GenerateNewToken() {
	user.Token = generateToken(user)
}

func (user *User) ChangeActiveStatus(newStatus bool) {
	user.IsActive = newStatus
}

func (user *User) ComparePassword(pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(pwd)) == nil
}

func generateToken(user *User) string {
	hash := md5.New()
	hash.Write([]byte(time.Now().String() + user.Email))
	return hex.EncodeToString(hash.Sum(nil))
}

func requiredIf(cond bool) validation.RuleFunc {
	return func (value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

func encryptString(s string) (string, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(enc), nil
}