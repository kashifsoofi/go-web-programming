package main

import (
	"net/http"
)

func process(w http.ResponseWriter, r *http.Request) {
	t, _ = t.ParseFiles("tmpl.html")
	content := `I asked: <i>"What's up?"</i>`
	t.Execute(w, content)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/process", process)
	server.ListenAndServe()
}
