package validate

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

var (
	errConfirmPassword = errors.New("Passwords don't match")
	errValidePassword  = errors.New("Invalid password")
)

// Pass ...
// Проверка пороля пользователя
func Pass(pass string) error {
	r, err := regexp.MatchString(`([A-Z1-9][a-z]+)`, pass)
	if err != nil {
		return err
	}

	if len(pass) < 8 || !r {
		return errValidePassword
	}

	return nil
}

// Confirm ...
// Проверка соответствия
func Confirm(first, last string) error {
	if first != last {
		return errConfirmPassword
	}

	return nil
}

// RequiredIf ...
// Обязательное поле
func RequiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

// AppIDValidate ...
// Провери id (Port) Приложения
func AppIDValidate(id string) error {
	return validation.Validate(
		id,
		validation.Required,
		validation.Match(regexp.MustCompile("^[0-9]{4}$")),
	)
}

// SystemName ...
// Проверить системное имя приложения
func SystemName(sys string) error {
	return validation.Validate(
		sys,
		validation.Required,
		validation.Match(regexp.MustCompile("^[a-z0-9_]{4,40}$")),
	)
}
