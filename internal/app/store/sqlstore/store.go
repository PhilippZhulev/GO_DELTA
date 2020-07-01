package sqlstore

import (
	"database/sql"

	"github.com/PhilippZhulev/delta/internal/app/store"
)

// Store ...
// Локализировать сущности
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
	testRepository *TestRepository
}

// New ...
// Создать новое хранилище
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

// Test ...
func (s *Store) Test() store.TestRepository {
	if s.testRepository != nil {
		return s.testRepository
	}

	s.testRepository = &TestRepository{
		store: s,
	}

	return s.testRepository
}
