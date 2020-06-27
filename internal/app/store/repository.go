package store

import (
	"database/sql"

	"github.com/PhilippZhulev/delta/internal/app/model"
)

// UserRepository ...
// Репозиторий для хранилища пользователя
type UserRepository interface {
	Create(*model.User) error
	FindByLogin(login string) (*model.User, error)
	Remove(id string) error
	GetAllUsers(l, o string) (*sql.Rows, error)
	Replace(*model.User) error
	ChangePassword(*model.User) error
}
