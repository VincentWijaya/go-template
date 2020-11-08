package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vincentwijaya/go-pkg/log"

	"github.com/go-chi/chi"
	"gopkg.in/gcfg.v1"
)

type Config struct {
	Server ServerConfig
	Log    LogConfig
}

type ServerConfig struct {
	Port        string
	Environment string
}

type LogConfig struct {
	LogPath string
	Level   string
}

const fileLocation = "/etc/serviceName/"
const devFileLocation = "files/etc/serviceName/"
const fileName = "serviceName.%s.yaml"
const devFileName = "serviceName.yaml"

const infoFile = "serviceName.info.log"
const errorFile = "serviceName.error.log"

func main() {
	//Read config
	var config Config
	location, fileName := getConfigLocation()
	err := gcfg.ReadFileInto(&config, location+fileName)
	if err != nil {
		log.Error("Failed to start service:", err)
		return
	}

	logConfig := log.LogConfig{
		StdoutFile: config.Log.LogPath + infoFile,
		StderrFile: config.Log.LogPath + errorFile,
		Level:      config.Log.Level,
	}
	log.InitLogger(config.Server.Environment, logConfig, []string{})

	// Repository

	// Usecase

	// Hanlder

	checker := systemCheck{
		pinger: map[string]Tester{},
	}

	httpRouter := chi.NewRouter()

	//CORS
	// httpRouter.Use(cors.Handler(cors.Options{
	// 	// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	ExposedHeaders:   []string{""},
	// 	AllowCredentials: false,
	// 	MaxAge:           300,
	// 	Debug:            true,
	// }))

	httpRouter.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("nothing here"))
	})
	httpRouter.Get("/ping", checker.ping)
	// httpRouter.Get("/health", checker.health)

	log.Infof("Service Started on:%v", config.Server.Port)
	err = http.ListenAndServe(config.Server.Port, httpRouter)
	if err != nil {
		log.Info("Failed serving Chi Dispatcher:", err)
		return
	}
	log.Info("Serving Chi Dispatcher on port:", config.Server.Port)
}

//-----------[ Pinger ]-----------------

type Tester interface {
	Ping() error
}

type systemCheck struct {
	pinger map[string]Tester
}

func (sys *systemCheck) ping(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

// DB Health pinger
// func (sys *systemCheck) health(w http.ResponseWriter, r *http.Request) {
// 	var str string
// 	for k, v := range sys.pinger {
// 		start := time.Now()
// 		status := "Success"
// 		message := "successful"
// 		if err := v.Ping(); err != nil {
// 			status = "Error"
// 			message = err.Error()
// 		}
// 		duration := time.Now().Sub(start).Milliseconds()
// 		str = fmt.Sprintf("%s%s | %s | %s | %dms\n", str, k, status, message, duration)
// 	}
// 	_, _ = w.Write([]byte(str))
// }

func getConfigLocation() (string, string) {
	env := os.Getenv("ENV")
	location := devFileLocation
	name := devFileName
	if env == "staging" || env == "production" || env == "development" {
		location = fileLocation
		name = fmt.Sprintf(fileName, env)
	}
	return location, name
}
