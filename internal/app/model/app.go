package model

import (
	"crypto/rand"
	"fmt"

	"github.com/PhilippZhulev/delta/internal/app/validate"
)

// App ...
// Структура приложения
type App struct {
	ID            int    `json:"id"`
	AppName       string `json:"app_name"`
	AppSystemName string `json:"app_system_name"`
	AppID         string `json:"app_id"`
	AppState      bool   `json:"app_state"`
	Rating        int    `json:"rating"`
	AppCategory   string `json:"app_category"`
	Token         string `json:"token"`
	AppDesc       string `json:"app_desc"`
	Avatar        string `json:"avatar"`
}

// AppLaunch ...
// Структура запущенных приложения
type AppLaunch struct {
	ID            int    `json:"id"`
	AppSystemName string `json:"app_system_name"`
	AppID         string `json:"app_id"`
	Pid           int    `json:"pid"`
}

// ValideID ...
// Проверка правильности ID
func (a App) ValideID(id string) error {
	return validate.AppIDValidate(id)
}

// ValideSystemName ...
// Проверка Имя приложения
func (a App) ValideSystemName(sys string) error {
	return validate.SystemName(sys)
}

// TokenGenerator ...
// Создать токен
func (a App) TokenGenerator() string {
	b := make([]byte, 25)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
