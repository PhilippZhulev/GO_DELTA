package model

import (
	"github.com/PhilippZhulev/delta/internal/app/validate"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID                       int    `json:"id"`
	Login                    string `json:"login"`
	Password                 string `json:"password,omitempty"`
	EncryptedPassword 			 string `json:"-"`
	СonfirmEncryptedPassword string `json:"confirm_password"`
	JobCode 					       string `json:"jobCode"`
	Email 						       string `json:"email"`
	Phone 									 string `json:"phone"`
	Name 										 string `json:"name"`
	UUID              			 string `json:"uuid"`
	Role 										 string `json:"role"`
}

// Validate ...
// Валидация при создание пользователя
func (u *User) Validate() error {

	err := validation.ValidateStruct(
		u, 
		validation.Field(&u.Login, validation.Required, validation.Length(4, 35)),
	) 

	if err != nil {
		return err
	} 

	return validate.Pass(u.EncryptedPassword)
}

// ValidatePassword ...
// Валидация пароля при изменении
func (u *User) ValidatePassword(first, last string) error {
	err := validate.Pass(first)
	if err != nil {
		return err
	}

	err = validate.Confirm(first, last)
	if err != nil {
		return err
	}

	return nil
}

// BeforeCreate ...
// Перед созданием пользователя
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}

	return nil
}

// Sanitize ...
// Очистка пароля
func (u *User) Sanitize() {
	u.Password = ""
}

// ComparePassword ...
// Хеширование
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
