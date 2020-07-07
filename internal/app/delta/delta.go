package delta

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/PhilippZhulev/delta/internal/app/confiiguration"
	"github.com/PhilippZhulev/delta/internal/app/store/sqlstore"
	"github.com/antonlindstrom/pgstore"
	_ "github.com/lib/pq" // ...
)

// Start ...
// Запустить сервер
func Start(config *confiiguration.Config) error {

	// Проверить соединение с Дб
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Создать хранилище
	store := sqlstore.New(db)

	// Создать хранилище сессий
	sessionStore, err := pgstore.NewPGStore(config.DatabaseURL, []byte(config.SessionKey))
	if err != nil {
		log.Fatal(err)
	}
	defer sessionStore.Close()
	defer sessionStore.StopCleanup(sessionStore.Cleanup(time.Minute * 5))

	// Очистка неактивных сессий
	go sessionStore.Cleanup(time.Hour * 2)

	// Создать сервер
	srv := newServer(store, sessionStore, config)
	server := &http.Server{Addr: config.BindAddr, Handler: srv}
	go func() {
			if err := server.ListenAndServe(); err != nil {
				log.Fatal(err)
			}
	}()

	// Отследить выключение сервера
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	// Изящное отключение сервера яерез shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
			return err
	}

	return nil
}

// Инициализировать новую базу данных
func newDB(dbURL string) (*sql.DB, error) {

	// Открыть соединение с db postgres
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
