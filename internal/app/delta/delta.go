package delta

import (
	"database/sql"
	"net/http"

	"github.com/PhilippZhulev/delta/internal/app/confiiguration"
	"github.com/PhilippZhulev/delta/internal/app/store/sqlstore"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq" // ...
)

// Start ...
// Запустить сервер
func Start(config *confiiguration.Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	
	srv := newServer(store, sessionStore, config)

	return http.ListenAndServe(config.BindAddr, srv)
}

// Инициализировать новую базу данных
func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
