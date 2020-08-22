package apiserver

import (
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

func Response(data []byte, w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err := w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
