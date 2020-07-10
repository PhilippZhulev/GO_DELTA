package handler

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/PhilippZhulev/delta/internal/app/helpers"
	"github.com/PhilippZhulev/delta/internal/app/model"
	"github.com/PhilippZhulev/delta/internal/app/store"
	"github.com/go-chi/chi"
)

//Статусы
var (
	startSuccess = "Application started..."
	stopSuccess  = "Application stoped..."
	appCreated   = "Application created"
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
	// Данные ответа
	type response struct {
		Pid int `json:"pid"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Установить root dir
		var (
			_, b, _, _ = runtime.Caller(0)
			root       = filepath.Join(filepath.Dir(b), "../../..")
		)
		// Запустить приложение через командную строку
		cmd := exec.Command(root+"/apps/test_app/app", "-port", chi.URLParam(r, "port"), "-name", "test_app_1")
		// Запуск cmd
		err := cmd.Start()
		if err != nil {
			ia.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}
		// Ответ
		//Отдает pid
		ia.respond.Done(w, r, http.StatusOK, &response{cmd.Process.Pid}, startSuccess)
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
		// Убить Процесс через cmd
		cmd := exec.Command("kill", chi.URLParam(r, "pid"))
		// Запустить коанду
		err := cmd.Run()
		if err != nil {
			ia.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}
		// Ответ
		ia.respond.Done(w, r, http.StatusOK, empty, stopSuccess)
	}
}
