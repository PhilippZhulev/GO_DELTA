package serve

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"net/url"
)

// Считывать флаги
// -port : порт для развочачивания
// -name : имя приложения
var (
	port = flag.String("port", "2101", "port")
	name = flag.String("name", "test_app", "application name")
)

// Request ...
// Структура запроса
type Request struct {
	Method string
	Param string
	Params map[string]string
	Session map[interface {}]interface {}
	Query url.Values
	Body string
	URL *url.URL
	Headers http.Header
	Context map[interface {}]interface {}
}

// Writer ...
// Структура для овета
type Writer struct {
	Data string
	Code int
	Msg string
}

// Method ...
// Структура для метода handler
type Method struct {
	Type string
}

// Error ...
// Метод заполнения ответа при ошибке
func (w *Writer) Error(code int, msg string) {
		w.Data = ""
		w.Code = code
		w.Msg = msg
}

// Send ...
//Метод заполнения ответа при успехе
func (w *Writer) Send(code int, data []byte, msg string) {
		w.Data = string(data)
		w.Code = code
		w.Msg = msg
}

// Use ...
// Промежуточный обработчик
func (r *Request) Use(w *Writer, callback func(r *Request, w *Writer)) {
	callback(r, w)
}

// Handle ...
// Обработчик запроса из delta
func (m Method) Handle(r Request, param string, callback func()) error {
	if (r.Method == m.Type) && (param == r.Param) {
		callback()
	}

	return nil
}

// Run ...
// Запустить tcp rpc сервер
func Run(delta interface{}) error {
	flag.Parse()

	p := *port
	n := *name

	err := rpc.Register(delta)
	if err != nil {
		return err
	}

	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":" + p)
	if err != nil {
		return err
	}

	fmt.Println("\033[1;33mService "+ n +" started...\033[0m")

	err = http.Serve(l, nil)
	if err != nil {
		return err
	}

	return nil
}