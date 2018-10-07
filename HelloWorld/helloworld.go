package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	var port = 8080
	println("Hello World in GO!!!")
	println("Running http server @ localhost:" + strconv.Itoa(port))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World in GO!!!")
	})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello again. You are @ %s\n", r.URL.Path)
	})
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
