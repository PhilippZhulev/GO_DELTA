package sqlstore

import (
	"database/sql"

	"github.com/PhilippZhulev/delta/internal/app/helpers"
	"github.com/PhilippZhulev/delta/internal/app/model"
	"github.com/PhilippZhulev/delta/internal/app/store"
)

// AppRepository ...
// Ссылка на хранилище
type AppRepository struct {
	hesh  helpers.Hesh
	store *Store
}

// Create ...
// Создать приложение
func (ar *AppRepository) Create(a *model.App) error {
	// Проверить формат ID
	if err := a.ValideID(a.AppID); err != nil {
		return err
	}
	// Проверить ситемное имя
	err := a.ValideSystemName(a.AppSystemName)
	if err != nil {
		return err
	}
	// Заполнить поля
	a.AppState = false
	a.Rating = 0
	a.Token = a.TokenGenerator()

	empty := "none"

	return ar.store.db.QueryRow(
		`
		INSERT INTO apps 
		(app_name, app_system_name, app_id, app_state, rating, app_category, token, , avatar) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id
		`,
		a.AppName,
		a.AppSystemName,
		a.AppID,
		a.AppState,
		a.Rating,
		a.AppCategory,
		a.Token,
		empty,
		empty,
	).Scan(&a.ID)
}

// GetAppToID ...
// Получить приложение по ID
func (ar *AppRepository) GetAppToID(a *model.App, al *model.AppLaunch, id string) error {
	// Проверить формат ID
	if err := a.ValideID(id); err != nil {
		return err
	}
	// Запрос в бд
	if err := ar.store.db.QueryRow(
		`
		SELECT app_system_name, app_id
		FROM apps 
		WHERE app_id = $1
		`,
		id,
	).Scan(
		&al.AppSystemName,
		&al.AppID,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}

	return nil
}

// LaunchApp ...
// Сохранить запущенное приложеник
func (ar *AppRepository) LaunchApp(al *model.AppLaunch) error {
	return ar.store.db.QueryRow(
		`
		INSERT INTO launch 
		(app_system_name, app_id, pid) 
		VALUES ($1, $2, $3) 
		RETURNING id
		`,
		al.AppSystemName,
		al.AppID,
		al.Pid,
	).Scan(&al.ID)
}

// GetLaunchApp ...
// Получить запущенное приложеник
func (ar *AppRepository) GetLaunchApp(al *model.AppLaunch, id string) error {
	// Запрос в бд
	if err := ar.store.db.QueryRow(
		`
		SELECT pid
		FROM launch 
		WHERE app_id = $1
		`,
		id,
	).Scan(&al.Pid); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}

	return nil
}

// RemoveLaunchApp ...
// Удалить запущенное приложеник
func (ar *AppRepository) RemoveLaunchApp(id string) error {
	// Запрос в бд
	_, err := ar.store.db.Exec(
		"DELETE FROM launch WHERE app_id = $1",
		id,
	)
	return err
}

// Change ...
// Изменить описание приложения
func (ar *AppRepository) Change(a *model.App, id string) error {
	// Проверить формат ID
	if err := a.ValideID(id); err != nil {
		return err
	}
	// Запрос в бд
	_, err := ar.store.db.Exec(
		`
		UPDATE apps
		SET app_name = $2, app_state = $3, app_desc = $4, app_category = $5
		WHERE app_id = $1
		`,
		id,
		a.AppName,
		a.AppState,
		a.AppDesc,
		a.AppCategory,
	)

	return err
}

// GetAllAppsAndFiltring ...
// Получить приложений по параметрам филтрации
func (ar *AppRepository) GetAllAppsAndFiltring(l, o, value string) (*sql.Rows, error) {
	return ar.store.db.Query(`
		SELECT id, app_system_name, app_name, app_category, rating 
		FROM apps  
		WHERE app_name LIKE '`+value+`%'
		LIMIT $1 
		OFFSET $2 * 2
	`, l, o)
}

// GetAllApps ...
// Получить приложения
func (ar *AppRepository) GetAllApps(l, o string) (*sql.Rows, error) {
	return ar.store.db.Query(`
		SELECT id, app_system_name, app_name, app_category, rating  
		FROM apps  
		LIMIT $1 
		OFFSET $2 * 2
	`, l, o)
}

// GetAppDataToID ...
// Получить приложение по ID
func (ar *AppRepository) GetAppDataToID(a *model.App, id string) error {
	// Проверить формат ID
	if err := a.ValideID(id); err != nil {
		return err
	}
	// Запрос в бд
	if err := ar.store.db.QueryRow(
		`
		SELECT app_name, app_system_name, app_id, app_state, app_category, rating, app_desc, avatar
		FROM apps 
		WHERE app_id = $1
		`,
		id,
	).Scan(
		&a.AppName,
		&a.AppSystemName,
		&a.AppID,
		&a.AppState,
		&a.AppCategory,
		&a.Rating,
		&a.AppDesc,
		&a.Avatar,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}

	return nil
}
