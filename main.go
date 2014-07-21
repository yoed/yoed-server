package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type yoedConfig struct {
	Listen string `json:"listen"`
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

func dispatch(handler, username, handle string) {

	resp, err := http.PostForm(handler, url.Values{
		"username": {username},
		"handle": {handle},
	})

	if err != nil {
		log.Printf("Error while dispatching message to %s: %s", handler, err)
		log.Printf("Remove handler %s", handler)
	} else {
		log.Printf("Handler %s status: %s", handler, resp.Status)
	}
}

func main() {

	config, err := loadConfig("./config.json")

	if err != nil {
		panic(fmt.Sprintf("failed loading config: %s", err))
	}

	handlers := make(map[string]map[string]bool, 0)

	router := mux.NewRouter()
	router.HandleFunc("/yo", func(w http.ResponseWriter, r *http.Request) {
		strHandles := r.FormValue("handles")
		callbackUrl := r.FormValue("callback_url")
		if strHandles == "" || callbackUrl == "" {
			errorMsg := "Handles and callback_url are mandatory"
			log.Printf("Error on subcribe: %s", errorMsg)
			http.Error(w, errorMsg, 400)
			return
		}

		handles := strings.Split(strHandles, ",")
		for _, handle := range handles {
			log.Printf("Subscribe %s: %s", handle, callbackUrl)
			if handlers[handle] == nil {
				handlers[handle] = make(map[string]bool, 0)
			}
			handlers[handle][callbackUrl] = true
		}
	})
	router.HandleFunc(`/yoed/{handle:[a-z0-9]+}`, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		handle := vars["handle"]
		username := r.FormValue("username")
		log.Printf("got a YO from %s on %s", username, handle)

		if 0 == len(handlers) || nil == handlers[handle] || 0 == len(handlers[handle]) {
			log.Printf("No handler registered for handle %s", handle)
			return
		}

		for handler, _ := range handlers[handle] {
			log.Printf("Dispatching to handler %s", handler)
			go dispatch(handler, username, handle)
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
