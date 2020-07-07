package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/PhilippZhulev/delta/internal/app/helpers"
	"github.com/PhilippZhulev/delta/internal/app/model"
	"github.com/PhilippZhulev/delta/internal/app/store"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/mitchellh/mapstructure"
	"github.com/sethvargo/go-password/password"
)

//Ошибки
var (
	errPagination = errors.New("Page does not exist")
)

//Статусы
var (
	userCreated = "User created"
	userRemoved = "User removed"
	sessionReceived = "Session received"
	userListReceived = "User list received"
	userReplace = "User replace"
	changePassword = "Password is changed"
	resetPassword = "Password is reset"
)


// InitUser ...
// Протокол аунтификации
type InitUser struct {
	respond *helpers.Respond
	store  store.Store
}

// HandleUserCreate ...
// Создать пользователя
func (iu *InitUser) HandleUserCreate(
	store store.Store,
) http.HandlerFunc {

	// Данные запроса
	type request struct {
		Name string `json:"name"`
		Login string `json:"login"`
		Password string `json:"password"`
		JobCode string `json:"jobCode"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	return func (w http.ResponseWriter, r *http.Request) {

		// Парсить запрос
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			iu.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Создать запись в store
		u := &model.User{
			Login:    req.Login,
			EncryptedPassword: req.Password,
			JobCode: req.JobCode,
			Name: req.Name,
			Email: req.Email,
			Phone: req.Phone,
		}
		
		// Записать пользователя в базу
		if err := store.User().Create(u); err != nil {
			iu.respond.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		// Если все ок отправить ответ
		iu.respond.Done(w, r, http.StatusCreated, u, userCreated)
	}
}

// HandleRemoveUser ...
// Удаление пользователя по id
func (iu *InitUser) HandleRemoveUser(store store.Store) http.HandlerFunc {

	// Данные ответа
	type respond struct {}

	return func (w http.ResponseWriter, r *http.Request) {

		// Параметр id
		id := chi.URLParam(r, "id")

		// Удалить пользователя из базы
		if err := store.User().Remove(id); err != nil {
			iu.respond.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		// Если все ок отправить ответ
		iu.respond.Done(w, r, http.StatusOK, &respond{}, userRemoved)
	}
}

// HandleUserSession ...
// Получить сессию пользователя
func (iu *InitUser) HandleUserSession(
	sesStore sessions.Store,
) http.HandlerFunc {

	// Данные ответа
	type respond struct {
		Name string `json:"name"`
		Login string `json:"login"`
		JobCode string `json:"jobCode"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Role string `json:"role"`
	}

	return func (w http.ResponseWriter, r *http.Request) {
		
		// Получить сессию
		session, err := sesStore.Get(r, iu.respond.GetUUID(r.Context()))
		if err != nil {
			iu.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Map to struct
		res := &respond{}
		err = mapstructure.Decode(session.Values, &res)
		if err != nil {
			iu.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Если все ок отправить ответ
		iu.respond.Done(w, r, http.StatusOK, res, sessionReceived)
	}
}

// HandleUserList ...
// Получить список пользователей
func (iu *InitUser) HandleUserList(store store.Store) http.HandlerFunc {

	// Структура элемента выборки
	type item struct {
		Name string `json:"name"`
		Login string `json:"login"`
		JobCode string `json:"jobCode"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Role string `json:"role"`
		UUID string `json:"uiid"`
		ID string `json:"id"`
	}

	// Данные ответа
	type respond struct {
		Size int `json:"size"`
		Page int `json:"page"`
		Result []item `json:"result"`
	}

	return func (w http.ResponseWriter, r *http.Request) {

		var (
			result []item
			im = &item{}
			l = chi.URLParam(r, "limit")
			o = chi.URLParam(r, "offset")
		) 

		// Запрос в бд
		usersRows, err := store.User().GetAllUsers(l, o);
		if err != nil {
			iu.respond.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		defer usersRows.Close()

		// Обработка выборки
		// генерация структуры
		// создане пегинации
		i := 0
		for usersRows.Next() {
			err := usersRows.Scan( 
				&im.ID, 
				&im.Login, 
				&im.JobCode, 
				&im.Email, 
				&im.Phone, 
				&im.Name, 
				&im.UUID, 
				&im.Role,
			)
			if err != nil {
				iu.respond.Error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			result = append(result, *im)
			i++
		}

		// Конверт параметров ссылки в int
		split, err := strconv.Atoi(l)
		page, err := strconv.Atoi(o)
		if err != nil {
			iu.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Заполнить ответ
		res := &respond{
			Result: result,
			Page: page,
			Size: split,
		}

		// Если все ок отправить ответ
		iu.respond.Done(w, r, http.StatusOK, res, userListReceived)
	}
}

//HandleUserReplace ...
//Изменение инфомации о пользователей
func (iu *InitUser) HandleUserReplace(store store.Store) http.HandlerFunc {

	// Структура запроса
	type request struct {
		Name string `json:"name"`
		Login string `json:"login"`
		JobCode string `json:"jobCode"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	// Структура ответа
	type response struct {
		Name string `json:"name"`
		JobCode string `json:"jobCode"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	return func (w http.ResponseWriter, r *http.Request)  {

		// Парсить запрос
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			iu.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Заполнить протокол
		u := &model.User{
			Login:   req.Login,
			JobCode: req.JobCode,
			Name: 	 req.Name,
			Email: 	 req.Email,
			Phone:   req.Phone,
		}

		// Перезапись полей в базе
		if err := store.User().Replace(u); err != nil {
			iu.respond.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		// Ответ
		res := &response{
			Name: 	 u.Name,
			JobCode: u.JobCode,
			Email: 	 u.Email,
			Phone: 	 u.Phone,
		}

		// Если все ок отправить ответ
		iu.respond.Done(w, r, http.StatusOK, res, userReplace)
	}
}

//HandleChangePassword ...
//Измененить пароль пользователя
func (iu *InitUser) HandleChangePassword(
	store store.Store, 
	sesStore sessions.Store,
) http.HandlerFunc {

	// Структура запроса
	type request struct {
		Password string  `json:"password"`
		EncryptedPassword string  `json:"new"`
		СonfirmEncryptedPassword string `json:"confirm"`
	}

	// Ответ
	type response struct {}

	return func (w http.ResponseWriter, r *http.Request)  {

		// Парсить запрос
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			iu.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Получить сессию
		session, err := sesStore.Get(r, iu.respond.GetUUID(r.Context()))
		if err != nil {
			iu.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Map to struct
		u := &model.User{}
		err = mapstructure.Decode(session.Values, &u)
		if err != nil {
			iu.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Заполнить данные пароля
		u.Password = req.СonfirmEncryptedPassword
		u.EncryptedPassword = req.EncryptedPassword
		u.СonfirmEncryptedPassword = req.СonfirmEncryptedPassword

		// Запросить изменение
		err = store.User().ChangePassword(u)
		if err != nil {
			iu.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}
		
		// Если все ок отправить ответ
		iu.respond.Done(w, r, http.StatusOK, &response{}, changePassword)
	}
}

//HandleResetPassword ...
// Сбросить пароль пользователя
func (iu *InitUser) HandleResetPassword(
	store store.Store, 
	sesStore sessions.Store,
) http.HandlerFunc {

	// Структура запроса
	type request struct {
		Login string  `json:"login"`
	}

	type response struct {
		Login string  `json:"login"`
		NewPassword string  `json:"newPassword"`
	}

	return func (w http.ResponseWriter, r *http.Request)  {

		// Парсить запрос
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			iu.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Сгенерировать пароль
		pass, err := password.Generate(8, 3, 0, false, true)
		if err != nil {
			iu.respond.Error(w, r, http.StatusBadGateway, err)
			return
		}

		// Заполнить данные пароля
		u := &model.User{}
		u.Login = req.Login
		u.EncryptedPassword = pass
		u.СonfirmEncryptedPassword = pass

		// Запросить изменение
		err = store.User().ChangePassword(u)
		if err != nil {
			iu.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Создать ответ`
		res := &response{
			Login: req.Login,
			NewPassword: pass,
		}
		
		// Если все ок отправить ответ
		iu.respond.Done(w, r, http.StatusOK, res, resetPassword)
	}
}