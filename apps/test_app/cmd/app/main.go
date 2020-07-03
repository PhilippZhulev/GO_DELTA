package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/delta/test_app/serve"
)

// Delta ...
//Инициализировать тип
type Delta string

// Тестовая структура для ответа
type respond struct {
		Name string `json:"name"`
		Session interface{} `json:"session"`
		Query string `json:"query"`
		Respond string `json:"respond"`
}

// Handler ...
// Пример запроса в delta
func (d *Delta) Handler(r serve.Request, w *serve.Writer) error {
	post := serve.Method{"POST"}

	// Создать обработчик
	post.Handle(r, "pinger", func() {

		// Заполняем структуру
		res := &respond{
			Name: "test", 
			Session: r.Session["login"],
			Query: r.Query.Get("name"),
			Respond: r.Params["data"] + "-pong",
		}

		// Структурв в JSON
		re, err := json.Marshal(res)
		if err != nil {

			// Error отправляет ответ c ошибкой в ядро delta
			// data при этом будет ровна null
			w.Error(http.StatusBadRequest, "Request Error")
			return
		}

		// Send отправляет ответ в ядро delta
		w.Send(http.StatusOK, re, "Request success")
	})

	return nil
}

func main() {
	//Запустить rpc обработчик delta
	delta := new(Delta)
	err := serve.Run(delta)
	if err != nil {
		log.Fatal(err)
	}
}