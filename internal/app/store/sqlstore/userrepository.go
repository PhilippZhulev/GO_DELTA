package sqlstore

import (
	"database/sql"

	"github.com/PhilippZhulev/delta/internal/app/helpers"
	"github.com/PhilippZhulev/delta/internal/app/model"
	"github.com/PhilippZhulev/delta/internal/app/store"
	"github.com/google/uuid"
)

// UserRepository ...
//Ссылка на хранилище
type UserRepository struct {
	hesh helpers.Hesh
	store *Store
}

// Create ...
//Создание пользователя в базе данных
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	// Заполнить допданные
	u.UUID = uuid.New().String()
	u.Role = "usr_default"
	u.EncryptedPassword = r.hesh.HashPassword(u.EncryptedPassword)

	// очистить пароль
	defer u.Sanitize()

	// Запрос в бд
	return r.store.db.QueryRow(
		`
		INSERT INTO users 
		(login_name, encrypted_password, jobcode, user_name, email, phone, uuid, role) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id
		`,
		u.Login,
		u.EncryptedPassword,
		u.JobCode,
		u.Name,
		u.Email,
		u.Phone,
		u.UUID,
		u.Role,
	).Scan(&u.ID)
}

// Remove ...
// Удаление пользователя
func (r *UserRepository) Remove(id string) error {
	
	// Запрос в бд
	_, err := r.store.db.Exec(
		"DELETE FROM users WHERE id = $1",
		id,
	);

	return err
}

// FindByLogin ...
// Поиск пользователя в базе данных
// Поск по логину
func (r *UserRepository) FindByLogin(login string) (*model.User, error) {
	u := &model.User{}

	// Запрос в бд
	if err := r.store.db.QueryRow(
		"SELECT * FROM users WHERE login_name = $1",
		login,
	).Scan(
		&u.ID,
		&u.Login,
		&u.EncryptedPassword,
		&u.JobCode,
		&u.Email,
		&u.Phone,
		&u.Name,
		&u.UUID,
		&u.Role,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// GetAllUsers ...
// Получить пользователей
func (r *UserRepository) GetAllUsers(l, o string) (*sql.Rows, error) {
	return r.store.db.Query(`
		SELECT id, login_name, jobcode, user_name, email, phone, uuid, role 
		FROM users  
		LIMIT $1 
		OFFSET $2 * 2
	`, l, o)
}

// Replace ...
// Изменить информацию пользователя
func (r *UserRepository) Replace(u *model.User) error {

	_, err := r.store.db.Exec(
		`
		UPDATE users
		SET jobcode = $2, user_name = $3, email = $4, phone = $5
		WHERE login_name = $1
		`,
		u.Login,
		u.JobCode,
		u.Name,
		u.Email,
		u.Phone,
	)

	// Запрос в бд
	return err
}

// ChangePassword ...
// Изменить пароль пользователя 
func (r *UserRepository) ChangePassword(u *model.User) error {

	err := u.ValidatePassword(u.EncryptedPassword, u.СonfirmEncryptedPassword)
	if err != nil {
		return err
	}

	// Заполнить доп данные
	u.EncryptedPassword = r.hesh.HashPassword(u.EncryptedPassword)

	// Запрос в бд
	_, err = r.store.db.Exec(
		`
		UPDATE users
		SET encrypted_password = $2
		WHERE login_name = $1
		`,
		u.Login,
		u.EncryptedPassword,
	)

	u.Sanitize()

	return err

}