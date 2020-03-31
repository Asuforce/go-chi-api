package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

type Handler func(http.ResponseWriter, *http.Request) (int, interface{}, error)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rv := recover(); rv != nil {
			debug.PrintStack()
			log.Printf("panic: %s", rv)
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}
	}()

	status, res, err := h(w, r)

	if err != nil {
		log.Printf("error: %s", err)
		RespondError(w, status, err)
	}
	respondJSON(w, status, res)
	return
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "   ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}
