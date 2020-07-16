package store

// Store ...
// Подключенные репозитории
type Store interface {
	User() UserRepository
	App() AppRepository
	Test() TestRepository
}
