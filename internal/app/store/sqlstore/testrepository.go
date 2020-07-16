package sqlstore

import (
	"database/sql"

	"github.com/PhilippZhulev/delta/internal/app/helpers"
)

// TestRepository ...
type TestRepository struct {
	hesh  helpers.Hesh
	store *Store
}

// GetTestRows ...
// Получить тестовые поля
func (t *TestRepository) GetTestRows() (*sql.Rows, error) {
	return t.store.db.Query(`SELECT * FROM test`)
}
