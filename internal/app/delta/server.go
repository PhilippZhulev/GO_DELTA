package delta

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/PhilippZhulev/delta/internal/app/confiiguration"
	"github.com/PhilippZhulev/delta/internal/app/handler"
	"github.com/PhilippZhulev/delta/internal/app/helpers"
	"github.com/PhilippZhulev/delta/internal/app/store"
)

// Протокол сервера
type server struct {
	router       *chi.Mux
	logger       *logrus.Logger
	store        store.Store
	tokenAuth    *jwtauth.JWTAuth
	sessionStore sessions.Store
	respond      *helpers.Respond
	
	user 				 handler.InitUser
	auth 				 handler.InitAuth
}

// Создать новый сервер
func newServer(
	store store.Store, 
	sessionStore sessions.Store, 
	config *confiiguration.Config,
) *server {

	// Заполнение протокола сервера
	s := &server{
		router:       chi.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		tokenAuth:    jwtauth.New("HS256", []byte(config.Salt), nil),
		sessionStore: sessionStore,
	}

	// Запуск middleWare
	s.configureMiddleware()

	return s
}

// Обслужить http
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// MiddleWare
func (s *server) configureMiddleware() {

	// middlewars
	s.router.Use(middleware.RequestID)
  s.router.Use(middleware.RealIP)
  s.router.Use(middleware.Logger)
  s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// Защищенный сегмент
	s.router.Group(func(router chi.Router) {

		// jwt
		router.Use(jwtauth.Verifier(s.tokenAuth))
		router.Use(s.authenticator)

		// Запросы пользователя
    router.Route("/api/v1/user", func(route chi.Router) {

			// Создание пользователя
			route.Post("/", s.user.HandleUserCreate(s.store))   
			
			//Изменить пользователя
			route.Put("/", s.user.HandleUserReplace(s.store))  

			// Удалить пользователя
			route.Delete("/{id}", s.user.HandleRemoveUser(s.store)) 

			// Создание пользователя
			route.Get("/", s.user.HandleUserSession(s.sessionStore))   

			// Получить списое пользоватей
			route.Get("/list/{limit}/{offset}", s.user.HandleUserList(s.store)) 

    })

	})

	// Открытый сегмент
	s.router.Group(func(router chi.Router) {
		// Запросы login
    router.Route("/api/v1/auth", func(route chi.Router) {

			// login
			route.Post("/login", s.auth.HandleLogin(s.store, s.sessionStore, s.tokenAuth))    

    })
	})
}

// Проверка jwt токена
// Промежуточное программное обеспеченеи
func (s server) authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		// Если какая либо ошибка
		if err != nil {
			s.respond.Error(w, r, http.StatusUnauthorized , err)
			return
		}

		// Если ошибка валидации
		if token == nil || !token.Valid {
			s.respond.Error(w, r, http.StatusUnauthorized , err)
			return
		}

		// Продолжить обработку
		next.ServeHTTP(w, r)
	})
}


