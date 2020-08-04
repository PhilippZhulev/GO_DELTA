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
	GetAllUsersAndFiltring(l, o, value string) (*sql.Rows, error)
	Replace(*model.User) error
	ChangePassword(*model.User) error
}

// AppRepository ...
// Репозиторий для приложений
type AppRepository interface {
	Create(*model.App) error
	GetAppToID(a *model.App, al *model.AppLaunch, id string) error
	LaunchApp(al *model.AppLaunch) error
	GetLaunchApp(al *model.AppLaunch, id string) error
	RemoveLaunchApp(id string) error
	Change(a *model.App, id string) error
	GetAllAppsAndFiltring(l, o, value string) (*sql.Rows, error)
	GetAllApps(l, o string) (*sql.Rows, error)
	GetAppDataToID(a *model.App, id string) error
}

// TestRepository ...
// Репозиторий для тестирования
type TestRepository interface {
	GetTestRows() (*sql.Rows, error)
}
