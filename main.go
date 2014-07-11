package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"github.com/gorilla/mux"
)

type yoedConfig struct {
	Listen   string `json:"listen"`
}

type yoedHandler interface {
	Handle(username string)
}

func loadConfig(configPath string) (*yoedConfig, error) {

	configFile, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	configJson, err := ioutil.ReadAll(configFile)

	if err != nil {
		return nil, err
	}

	config := &yoedConfig{}

	if err := json.Unmarshal(configJson, config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {

	config, err := loadConfig("./config.json")

	if err != nil {
		panic(fmt.Sprintf("failed loading config: %s", err))
	}

	handlers := make([]string, 0)

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		callbackUrl := r.FormValue("callback_url")
		log.Printf("subscribe %s", callbackUrl)
		handlers = append(handlers, callbackUrl)
	})
	router.HandleFunc(`/yoed/{handle:[a-z0-9]+}`, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		handle := vars["handle"]
		username := r.FormValue("username")
		log.Printf("got a YO from %s on %s", username, handle)

		for _, handler := range handlers {
			log.Printf("Dispatch to handler %s", handler)
			http.PostForm(handler, url.Values{"username":{username}})
		}
	})

	server := http.Server{
		Addr:    config.Listen,
		Handler: router,
	}

	log.Printf("Listening...")

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}

}