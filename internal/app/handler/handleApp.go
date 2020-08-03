package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/PhilippZhulev/delta/internal/app/helpers"
	"github.com/PhilippZhulev/delta/internal/app/model"
	"github.com/PhilippZhulev/delta/internal/app/store"
	"github.com/go-chi/chi"
)

//Статусы
var (
	startSuccess    = "Application started..."
	stopSuccess     = "Application stoped..."
	appCreated      = "Application created"
	appChanged      = "Application change"
	appListReceived = "Application list received"
)

// InitApp ...
// Протокол аунтификации
type InitApp struct {
	respond *helpers.Respond
	store   store.Store
}

// CreateApp ...
// Создать приложение
func (ia InitApp) CreateApp(store store.Store) http.HandlerFunc {
	// Данные запроса
	type request struct {
		AppSystemName string `json:"appSystemName"`
		AppName       string `json:"appName"`
		AppCategory   string `json:"appCategory"`
		AppID         string `json:"appId"`
	}
	//Данные ответа
	type response struct {
		ID            int    `json:"id"`
		Tonen         string `json:"token"`
		AppSystemName string `json:"appSystemName"`
		AppID         string `json:"appId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Парсить запрос
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			ia.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}
		// Создать запись в store
		a := &model.App{
			AppSystemName: req.AppSystemName,
			AppName:       req.AppName,
			AppCategory:   req.AppCategory,
			AppID:         req.AppID,
		}
		// Записать приложение в базу
		if err := store.App().Create(a); err != nil {
			ia.respond.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		// Если все успешно отправить ответ
		ia.respond.Done(w, r, http.StatusOK, &response{a.ID, a.Token, a.AppSystemName, a.AppID}, appCreated)
	}
}

// RunApplication ...
// Активировать приложение
func (ia InitApp) RunApplication(store store.Store) http.HandlerFunc {
	//Данные ответа
	type response struct {
		ID            int    `json:"id"`
		AppSystemName string `json:"appSystemName"`
		AppID         string `json:"appId"`
		Pid           int    `json:"pid"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Установить root dir
		var (
			_, b, _, _ = runtime.Caller(0)
			root       = filepath.Join(filepath.Dir(b), "../../..")
			port       = chi.URLParam(r, "port")
		)

		// Удалить запись о запуске
		// если приложениу уже было запущенно
		_ = store.App().RemoveLaunchApp(port)

		// Создать запись в store
		al := &model.AppLaunch{}
		a := &model.App{}
		// Записать приложение в базу
		if err := store.App().GetAppToID(a, al, port); err != nil {
			ia.respond.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		// Запустить приложение через командную строку
		cmd := exec.Command(root+"/apps/"+al.AppSystemName+"/app", "-port", port, "-name", al.AppSystemName)
		if err := cmd.Start(); err != nil {
			ia.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		al.Pid = cmd.Process.Pid
		// Записать приложение в базу
		if err := store.App().LaunchApp(al); err != nil {
			ia.respond.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		// Ответ
		//Отдает pid
		res := &response{al.ID, al.AppSystemName, al.AppID, al.Pid}
		ia.respond.Done(w, r, http.StatusOK, res, startSuccess)
	}
}

// StopApplication ...
// Активировать приложение
func (ia InitApp) StopApplication(store store.Store) http.HandlerFunc {
	// Данные ответа
	type response struct {
		State bool `json:"state"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// port
		port := chi.URLParam(r, "port")
		al := &model.AppLaunch{}
		// Получить pid
		if err := store.App().GetLaunchApp(al, port); err != nil {
			ia.respond.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		// Убить Процесс через cmd
		cmd := exec.Command("kill", strconv.Itoa(al.Pid))
		// Запустить коанду
		if err := cmd.Run(); err != nil {
			ia.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}
		// Удалить из launch
		if err := store.App().RemoveLaunchApp(port); err != nil {
			ia.respond.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		// Ответ
		ia.respond.Done(w, r, http.StatusOK, empty, stopSuccess)
	}
}

// ChangeApp ...
// Изменить даные приложения приложение
func (ia InitApp) ChangeApp(store store.Store) http.HandlerFunc {
	// Данные запроса
	type request struct {
		ID          string `json:"id"`
		AppName     string `json:"name"`
		AppCategory string `json:"category"`
		AppDesc     string `json:"desc"`
		AppState    bool   `json:"state"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Парсить запрос
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			ia.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}
		// Создать запись в store
		a := &model.App{
			AppName:     req.AppName,
			AppCategory: req.AppCategory,
			AppDesc:     req.AppDesc,
			AppState:    req.AppState,
		}
		// Записать приложение в базу
		if err := store.App().Change(a, req.ID); err != nil {
			ia.respond.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		// Ответ
		ia.respond.Done(w, r, http.StatusOK, &req, appChanged)
	}
}

// GetAppList ...
// Получить список приложений
func (ia InitApp) GetAppList(store store.Store) http.HandlerFunc {
	// Данные элемента возврата
	type item struct {
		AppID         string `json:"id"`
		AppSystemName string `json:"systemName"`
		AppName       string `json:"name"`
		AppCategory   string `json:"category"`
		AppRating     int    `json:"rating"`
	}
	// Тело запроса
	// если POST
	type request struct {
		Value string `json:"value"`
	}
	// Данные ответа
	type respond struct {
		Size   int    `json:"size"`
		Page   int    `json:"page"`
		Result []item `json:"result"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Основные переменные
		var (
			result  []item
			im      = &item{}
			l       = chi.URLParam(r, "limit")
			o       = chi.URLParam(r, "offset")
			appRows *sql.Rows
			err     error
		)
		// Запрос в бд
		if r.Method == "POST" {
			// Парсить запрос
			req := &request{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				ia.respond.Error(w, r, http.StatusBadRequest, err)
				return
			}
			// Запрос + фильтрация + пегинация
			appRows, err = store.App().GetAllAppsAndFiltring(l, o, req.Value)
		} else {
			// Запрос + пегинация
			appRows, err = store.App().GetAllApps(l, o)
		}

		if err != nil {
			ia.respond.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		defer appRows.Close()
		// Обработка выборки
		// генерация структуры
		// создане пегинации
		i := 0
		for appRows.Next() {
			err := appRows.Scan(
				&im.AppID,
				&im.AppSystemName,
				&im.AppName,
				&im.AppCategory,
				&im.AppRating,
			)
			if err != nil {
				ia.respond.Error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			result = append(result, *im)
			i++
		}
		// Конверт параметров ссылки в int
		split, err := strconv.Atoi(l)
		page, err := strconv.Atoi(o)
		if err != nil {
			ia.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}
		// Заполнить ответ
		res := &respond{
			Result: result,
			Page:   page,
			Size:   split,
		}
		// Если все ок отправить ответ
		ia.respond.Done(w, r, http.StatusOK, res, appListReceived)
	}
}
