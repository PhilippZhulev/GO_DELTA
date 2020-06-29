package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/PhilippZhulev/delta/internal/app/confiiguration"
	"github.com/PhilippZhulev/delta/internal/app/delta"
	"github.com/sevlyar/go-daemon"
)

// Путь к config
var (
	configPath string
)

// Инициализировать
// ----
// Заполнение config
func init() {
	flag.StringVar(&configPath, "config-path", "configs/delta.toml", "path to config file")
}

// Инициализировать сервер
// ----
// Создать конфиг
// Начать инициализацию сервера с парам.
// конфига
func serve() {
	flag.Parse()

	config := confiiguration.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	srv := delta.Start(config); 

	err = srv.Shutdown(nil)
	if err != nil {
			log.Fatal(err)
	}
}

// main
// ----
// Запуск демона
// Запуск ининциализации сервера
func main() {
	fmt.Println("\033[1;33mDelta api started...\033[0m")

	deamon := flag.Bool("deamon", false, "a bool")
	flag.Parse()
	
	if *deamon == true {
		cntxt := &daemon.Context{
			PidFileName: "sample.pid",
			PidFilePerm: 0644,
			LogFileName: "demon.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
			Args:        []string{"[go-daemon sample]"},
		}

		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			return
		}
		defer cntxt.Release()

		log.Print("- - - - - - - - - - - - - - -")
		log.Print("daemon started")
	}

	serve()
}
