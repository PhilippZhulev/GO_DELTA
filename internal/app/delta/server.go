package delta

import (
	"errors"
	"net/http"
	"time"

	"github.com/PhilippZhulev/delta/internal/app/confiiguration"
	"github.com/PhilippZhulev/delta/internal/app/handler"
	"github.com/PhilippZhulev/delta/internal/app/helpers"
	"github.com/PhilippZhulev/delta/internal/app/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

//Ошибки
var (
	errNoSession = errors.New("session does not exist")
)

// Протокол сервера
type server struct {
	router       *chi.Mux
	logger       *logrus.Logger
	store        store.Store
	tokenAuth    *jwtauth.JWTAuth
	sessionStore sessions.Store
	respond      *helpers.Respond

	user     handler.InitUser
	auth     handler.InitAuth
	dispatch handler.InitDispatch
	app      handler.InitApp
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
		router.Use(s.authenticator(s.sessionStore))
		// Диспатчер приложений
		router.Route("/dispatch", func(route chi.Router) {
			route.Post("/{port}/{param}/*", s.dispatch.HandleDispatch(s.sessionStore))
			route.Put("/{port}/{param}/*", s.dispatch.HandleDispatch(s.sessionStore))
			route.Get("/{port}/{param}/*", s.dispatch.HandleDispatch(s.sessionStore))
			route.Delete("/{port}/{param}/*", s.dispatch.HandleDispatch(s.sessionStore))
			route.Options("/{port}/{param}/*", s.dispatch.HandleDispatch(s.sessionStore))
		})
		// Управление приложением
		router.Route("/api/v1/app", func(route chi.Router) {
			route.Options("/run/{port}", s.app.RunApplication(s.store))
			route.Options("/stop/{port}", s.app.StopApplication(s.store))
			route.Post("/", s.app.CreateApp(s.store))
			route.Put("/", s.app.ChangeApp(s.store))
		})
		// Запросы пользователя
		router.Route("/api/v1/user", func(route chi.Router) {
			// Создание пользователя
			route.Post("/", s.user.HandleUserCreate(s.store))
			//Изменить пользователя
			route.Put("/", s.user.HandleUserReplace(s.store))
			//Изменить пароль пользователя
			route.Put("/password", s.user.HandleChangePassword(s.store, s.sessionStore))
			//Сбросить пользователя
			route.Patch("/password", s.user.HandleResetPassword(s.store, s.sessionStore))
			// Удалить пользователя
			route.Delete("/{id}", s.user.HandleRemoveUser(s.store))
			// Создание пользователя
			route.Get("/", s.user.HandleUserSession(s.sessionStore))
			// Получить список пользоватей
			route.Get("/list/{limit}/{offset}", s.user.HandleUserList(s.store))
			// Получить список пользоватей (Фильтрация)
			route.Post("/list/{limit}/{offset}", s.user.HandleUserList(s.store))
		})
		// Запросы auth
		router.Route("/api/v1/auth", func(route chi.Router) {
			// logout
			route.Get("/logout", s.auth.HandleLogout(s.store, s.sessionStore, s.tokenAuth))
			// get active sessions
			route.Post("/session", s.auth.CheckActiveSession(s.store, s.sessionStore))
		})

	})
	// Открытый сегмент
	s.router.Group(func(router chi.Router) {
		// Запросы auth
		router.Route("/api/v1/auth/login", func(route chi.Router) {
			// login
			route.Post("/", s.auth.HandleLogin(s.store, s.sessionStore, s.tokenAuth))
		})
	})
}

// Проверка jwt токена
// Промежуточное программное обеспеченеи
func (s server) authenticator(sesStore sessions.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получить токен
			token, cl, err := jwtauth.FromContext(r.Context())
			if err != nil {
				s.respond.Error(w, r, http.StatusUnauthorized, err)
				return
			}
			// Получить сессию
			session, _ := sesStore.Get(r, cl["uuid"].(string))
			// Если какая либо ошибка
			if err != nil {
				s.respond.ClearSession(session, w, r)
				s.respond.Error(w, r, http.StatusUnauthorized, err)
				return
			}
			// Если ошибка валидации
			if token == nil || !token.Valid {
				s.respond.ClearSession(session, w, r)
				s.respond.Error(w, r, http.StatusUnauthorized, err)
				return
			}
			// Проверить сессию
			uuid := session.Values["uuid"]
			if uuid == nil {
				s.respond.Error(w, r, http.StatusUnauthorized, errNoSession)
				return
			}
			// Продолжить обработку
			next.ServeHTTP(w, r)
		})
	}
}
