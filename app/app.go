package app

import (
	"log"
	"net/http"

	"github.com/Prajithp/gchat-notifier/app/handler"
	"github.com/Prajithp/gchat-notifier/config"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	Config *config.Config
}

func (a *App) Initialize(config *config.Config) {
	a.Router = mux.NewRouter()
	a.Config = config
	a.setRouters()
}

func (a *App) setRouters() {
	a.Post("/send/{channel}", a.handleRequest(handler.Notification))
}

func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(config *config.Config, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.Config, w, r)
	}
}
