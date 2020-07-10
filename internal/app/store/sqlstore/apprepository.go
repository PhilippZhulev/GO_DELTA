package sqlstore

import (
	"github.com/PhilippZhulev/delta/internal/app/helpers"
	"github.com/PhilippZhulev/delta/internal/app/model"
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
	err := a.ValideID(a.AppID)
	if err != nil {
		return err
	}
	// Проверить ситемное имя
	err = a.ValideSystemName(a.AppSystemName)
	if err != nil {
		return err
	}
	// Заполнить поля
	a.AppState = false
	a.Rating = 0
	a.Token = a.TokenGenerator()

	return ar.store.db.QueryRow(
		`
		INSERT INTO apps 
		(app_name, app_system_name, app_id, app_state, rating, app_category, token) 
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
	).Scan(&a.ID)
}
