package serve

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"net/url"
)

var (
	port = flag.String("port", "2101", "port")
	name = flag.String("name", "test_app", "application name")
)

// Request ...
type Request struct {
	Method string
	Param string
	Params map[string]string
	Session map[interface {}]interface {}
	Query url.Values
	Body string
	URL *url.URL
	Headers http.Header
}

// Writer ...
type Writer struct {
	Data string
	Code int
	Msg string
}

// Method ...
type Method struct {
	Type string
}

// Error ...
func (w *Writer) Error(code int, msg string) {
		w.Data = ""
		w.Code = code
		w.Msg = msg
}

// Send ...
func (w *Writer) Send(code int, data []byte, msg string) {
		w.Data = string(data)
		w.Code = code
		w.Msg = msg
}

// Handle ...
func (m Method) Handle(r Request, param string, callback func()) error {
	if (r.Method == m.Type) && (param == r.Param) {
		callback()
	}

	return nil
}

// Run ...
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