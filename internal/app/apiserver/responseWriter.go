package apiserver

import (
	"encoding/json"
	"log"
	"net/http"
)

func NewResponseWriter(json []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(json)
	if err != nil {
		log.Println(err)
	}
}

func Reponse(body []byte, w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(err.Error())
	}

	_, err := w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
