package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type UnspecifiedMessageError struct{}

func (m *UnspecifiedMessageError) Error() string {
	return "Unspecified Message"
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debug(string(body))

	response, err := NewWebhookResponse(body)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sender, err := Sender(nil), &UnspecifiedMessageError{};


	if response.Config != nil {
		config, _ := json.Marshal(response.Config)
		sender, err = senderDispatch("CONFIG", response, []byte(config))
	} else if response.Message != "" {
		sender, err = senderDispatch("MESSAGE", response, []byte(response.Message))
	}
	
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profiles := viper.GetStringSlice("profiles")
	for _, profile := range profiles {
		if profile == "email" {
			if err := sender.SendEmail(); err != nil {
				log.Error(err)
				return
			}
		}
		if profile == "slack" {
			if err := sender.SendSlack(); err != nil {
				log.Error(err)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	addr := net.JoinHostPort(viper.GetString("listen_host"), viper.GetString("listen_port"))
	log.Infof("Server listening on %s%s", addr, viper.GetString("webhook_path"))
	http.ListenAndServe(addr, http.HandlerFunc(handler))
}
