package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Respond ...
// Протокол хелперов
type Respond struct {
	Data interface{} `json:"data"`
	Msg string `json:"msg"`
}

// Error ...
// Ошибка запроса
func (h *Respond) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	var ed map[string]string
	h.Done(w, r, code, ed, err.Error())
}

// Done ...
// Ответ сервера
func (h *Respond) Done(w http.ResponseWriter, r *http.Request, code int, data interface{}, mess string) {
	result := &Respond{
		Data: data,
		Msg: mess,
	}
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(result)
	}
}

// ParseDone ...
// Ответ сервера c массивом
func (h *Respond) ParseDone(w http.ResponseWriter, r *http.Request, code int, data []string, mess string) {
	result := `{"data": ${data}, "msg": "${mess}" }`
	result = strings.ReplaceAll(result, "${data}", strings.Join(data, ""))
	result = strings.ReplaceAll(result, "${mess}", mess)

	w.WriteHeader(code)
  w.Write([]byte(fmt.Sprintf(result)))
}