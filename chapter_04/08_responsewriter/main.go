package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeBody(w http.ResponseWriter, r *http.Request) {
	str := `
	<html>
		<head>
			<title>Go Web Programming</title>
		</head>
		<body>
			<h1>Hello World</h1>
		</body>
	</html>`
	w.Write([]byte(str))
}

func writeHeader(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	fmt.Fprintf(w, "No such service, try next door")
}

func header(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "http://google.com")
	w.WriteHeader(http.StatusFound)
}

type Post struct {
	User    string
	Threads []string
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	post := Post{
		User:    "Kashif Jamal",
		Threads: []string{"first", "second", "third"},
	}
	json, _ := json.Marshal(post)
	w.Write(json)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/write", writeBody)
	http.HandleFunc("/writeHeader", writeHeader)
	http.HandleFunc("/redirect", header)
	http.HandleFunc("/post", post)
	server.ListenAndServe()
}
