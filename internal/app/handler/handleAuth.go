package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/PhilippZhulev/delta/internal/app/helpers"
	"github.com/PhilippZhulev/delta/internal/app/store"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/sessions"
)

//Ошибки
var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

//Статусы
var (
	successLogin = "Login success"
	successLogout = "Logout success"
)

// InitAuth ...
// Протокол аунтификации
type InitAuth struct {
	respond *helpers.Respond
	hesh helpers.Hesh
}

// HandleLogin ...
// Логин
func (ia InitAuth) HandleLogin(
	store store.Store, 
	sesStore sessions.Store,
	tokenAuth *jwtauth.JWTAuth,
) http.HandlerFunc {

	// Данные запоса
	type request struct {
		AuthData string `json:"authData"`
	}

	// Данные ответа
	type response struct {
		Token string `json:"token"`
	}

	return func (w http.ResponseWriter, r *http.Request) {

		// Распарсить запрос
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			ia.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Декодировать auth данные
		decoded, err := base64.StdEncoding.DecodeString(req.AuthData)
		if err != nil {
				ia.respond.Error(w, r, http.StatusBadRequest, err)
				return
		}
		result := strings.Split(string(decoded), ":")

		// Искать юзера
		u, err := store.User().FindByLogin(result[0])
		if err != nil {
			ia.respond.Error(w, r, http.StatusBadRequest, errIncorrectEmailOrPassword)
			return
		}

		// Проверить соответстие пароля
		if !ia.hesh.CheckPasswordHash(result[1], u.EncryptedPassword) {
			ia.respond.Error(w, r, http.StatusBadRequest, errIncorrectEmailOrPassword)
			return 
		}

		// Получить сессию
		session, err := sesStore.Get(r, "delta_session")
		if err != nil {
			ia.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Заполнить сессию
		session.Values["id"] = u.ID
		session.Values["login"] = u.Login
		session.Values["name"] = u.Name
		session.Values["jobcode"] = u.JobCode
		session.Values["email"] = u.Email
		session.Values["phone"] = u.Phone
		session.Values["role"] = u.Role
		session.Values["uuid"] = u.UUID

		// Сохранить сессию
		err = session.Save(r, w)
		if err != nil {
			ia.respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Сгенерировать токен
		claims := jwt.MapClaims{"user_id": u.ID, "uuid": u.UUID}
		jwtauth.SetExpiryIn(claims, 1000 * time.Second)
		_, tokenString, err := tokenAuth.Encode(claims)
		if err != nil {
			ia.respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Если все успешно дать ответ
		// отдает токен
		ia.respond.Done(w, r, http.StatusOK, &response{tokenString}, successLogin)
	}
}

// HandleLogout ...
// Выйти из системы
func (ia InitAuth) HandleLogout(
	store store.Store, 
	sesStore sessions.Store,
	tokenAuth *jwtauth.JWTAuth,
) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {

		// Получить сессию
		session, err := sesStore.Get(r, "delta_session")
		if err != nil {
			ia.respond.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		// Удалить сессию
		session.Options.MaxAge = -1
		err = session.Save(r, w)
		if err != nil {
			ia.respond.Error(w, r, http.StatusInternalServerError, err)
		}

		// Пустой data
		var empty []string

		// Если все успешно дать ответ
		ia.respond.Done(w, r, http.StatusOK, empty, successLogout)
	}
}