package handler

import (
	"io/ioutil"
	"net/http"
	"net/rpc"
	"net/url"

	"github.com/PhilippZhulev/delta/internal/app/helpers"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
)

// InitDispatch ...
// Протокол аунтификации
type InitDispatch struct {
	respond *helpers.Respond
	hesh helpers.Hesh
}

// HandleDispatch ...
// Запрос приложение
func (idch *InitDispatch) HandleDispatch(
	sesStore sessions.Store,
) http.HandlerFunc {

	// Ответ
	type Result struct {
		Data string
		Code int
		Msg string
	}

	// Данные ответа
	type respond struct {
		Result Result
	}

	// Args ...
	// urlParams
	type Args struct {
		Method string
		Param string
		Session map[interface {}]interface {}
		Params map[string]string
		Query url.Values
		Body string
		URL *url.URL
		Headers http.Header
		Context map[interface {}]interface {}
	}

	return func (w http.ResponseWriter, r *http.Request) {

		// Параметры ссылки
		// -port (port - id) sub.app
		// -param параметор url для индификации в sub.app
		port := chi.URLParam(r, "port")
		p := chi.URLParam(r, "param")

		// Открть соединение rpc с дочерним сервером
		// В качестве порта взять параметр port
		client, err := rpc.DialHTTP("tcp", ":" + port)
		if err != nil {
			idch.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Получить контекст
		// Получить все параметры ссылки
		ctx := chi.RouteContext(r.Context())
		prm := make(map[string]string)
		for i := 0; i < len(ctx.URLParams.Keys); i++ {
			prm[ctx.URLParams.Keys[i]] = ctx.URLParams.Values[i]
		} 

		// Получить данные из body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			idch.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		// Получить сессию
		session, err := sesStore.Get(r, idch.respond.GetUUID(r.Context()))
		if err != nil {
			idch.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Созать мапу для контекста
		con := make(map[interface{}]interface{})

		// Аргументы для передачи sub.app
		// передает рание получинные параметры запроса
		args := &Args{
			Method: r.Method,
			Session: session.Values,
			Param: p,
			Params: prm,
			Query: r.URL.Query(),
			Body: string(body),
			URL: r.URL,
			Headers: r.Header,
			Context: con,
		}

		// Отправить параметры в handler sub.app
		result := &Result{}
		err = client.Call("Delta.Handler", args, &result)
		if err != nil {
			idch.respond.Error(w, r, http.StatusBadRequest, err)
			return
		}

		// Если все успешно дать ответ
		// Если handler sub.app возвращает пустой data
		// отправить null 
		if len(result.Data) > 0 {
			idch.respond.ParseDone(w, r, result.Code, result.Data, result.Msg)
		} else  {
			var empty []string
			idch.respond.Done(w, r, result.Code, empty, result.Msg)
		}
	}
}

