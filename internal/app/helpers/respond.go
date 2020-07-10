package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/jwtauth"
	"github.com/gorilla/sessions"
)

// Respond ...
// Протокол хелперов
type Respond struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
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
		Msg:  mess,
	}
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(result)
	}
}

// ParseDone ...
// Ответ сервера c массивом
func (h *Respond) ParseDone(w http.ResponseWriter, r *http.Request, code int, data string, mess string) {
	result := `{"data": ${data}, "msg": "${mess}" }`
	result = strings.ReplaceAll(result, "${data}", data)
	result = strings.ReplaceAll(result, "${mess}", mess)

	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(result)))
}

// ClearSession ...
// Очистка сессии
func (h *Respond) ClearSession(s *sessions.Session, w http.ResponseWriter, r *http.Request) {
	// Удалить сессию
	s.Options.MaxAge = -1
	_ = s.Save(r, w)
}

// GetUUID ...
// Получить uuid из токена
func (h *Respond) GetUUID(ctx context.Context) string {
	_, cl, _ := jwtauth.FromContext(ctx)
	return cl["uuid"].(string)
}
