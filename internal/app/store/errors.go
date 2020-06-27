package store

import "errors"

var (
	// ErrRecordNotFound ...
	// Ошибка получения записи
	ErrRecordNotFound = errors.New("record not found")
)
