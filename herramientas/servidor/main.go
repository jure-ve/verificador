package main

import (
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/lento/", func(w http.ResponseWriter, r *http.Request) { time.Sleep(time.Second); w.Write([]byte("ok")) })
	http.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(503) })
	http.HandleFunc("/colgado", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(10 * time.Second):
		case <-r.Context().Done():
		}
	})
	http.ListenAndServe("127.0.0.1:8099", nil)
}
