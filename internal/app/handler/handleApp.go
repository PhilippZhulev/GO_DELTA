package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/PhilippZhulev/delta/internal/app/helpers"
	"github.com/PhilippZhulev/delta/internal/app/model"
	"github.com/PhilippZhulev/delta/internal/app/store"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

//Статусы
var (
	startSuccess    = "Application started..."
	stopSuccess     = "Application stoped..."
	appCreated      = "Application created"
	appChanged      = "Application change"
	appListReceived = "Application list received"
	appReceived     = "Application received"
	upload          = "File upload"
)

// InitApp ...
// Протокол аунтификации
type InitApp struct {
	hesh    helpers.Hesh
	respond *helpers.Respond
	store   store.Store
}

// HandleCreateApp ...
// Создать приложение
func (ia InitApp) HandleCreateApp(store store.Store) http.HandlerFunc {
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

// HandleRunApplication ...
// Активировать приложение
func (ia InitApp) HandleRunApplication(store store.Store) http.HandlerFunc {
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

// HandleStopApplication ...
// Активировать приложение
func (ia InitApp) HandleStopApplication(store store.Store) http.HandlerFunc {
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

// HandleChangeApp ...
// Изменить даные приложения приложение
func (ia InitApp) HandleChangeApp(store store.Store) http.HandlerFunc {
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

// HandleGetAppList ...
// Получить список приложений
func (ia InitApp) HandleGetAppList(store store.Store) http.HandlerFunc {
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

// HandleGetApp ...
// Получить данные приложения
func (ia InitApp) HandleGetApp(store store.Store) http.HandlerFunc {
	// Данные элемента возврата
	type respond struct {
		AppID         string `json:"id"`
		AppSystemName string `json:"systemName"`
		AppName       string `json:"name"`
		AppCategory   string `json:"category"`
		AppRating     int    `json:"rating"`
		AppState      bool   `json:"state"`
		AppDesc       string `json:"desc"`
		Avatar        string `json:"avatar"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		a := &model.App{}
		// Записать приложение в базу
		if err := store.App().GetAppDataToID(a, id); err != nil {
			ia.respond.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		// Заполнить ответ
		res := &respond{
			AppID:         a.AppID,
			AppSystemName: a.AppSystemName,
			AppName:       a.AppName,
			AppCategory:   a.AppCategory,
			AppRating:     a.Rating,
			AppState:      a.AppState,
			AppDesc:       a.AppDesc,
			Avatar:        a.Avatar,
		}
		// Если все ок отправить ответ
		ia.respond.Done(w, r, http.StatusOK, res, appReceived)
	}
}

// HandleAddPreview ...
// Добавить превью
func (ia InitApp) HandleAddPreview(store store.Store) http.HandlerFunc {
	//Ответ
	type respond struct {
		Size int64  `json:"size"`
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Получить файл из формы
		file, handler, err := r.FormFile("avatar")
		if err != nil {
			ia.respond.Done(w, r, http.StatusOK, r.Body, appReceived)
			return
		}
		defer file.Close()
		// Создать новый файл
		name := uuid.New().String() + strings.Split(handler.Filename, ".")[1]
		f, err := os.OpenFile("./upload/"+name, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			ia.respond.Done(w, r, http.StatusOK, r.Body, appReceived)
			return
		}
		defer f.Close()
		// Копировать содержимое файла
		io.Copy(f, file)
		// Собираю информацию о файле и отправляю в ответе
		fi, err := f.Stat()
		res := respond{fi.Size(), name}
		ia.respond.Done(w, r, http.StatusOK, res, upload)
	}
}
