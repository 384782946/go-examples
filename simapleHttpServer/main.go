package main

import (
	"fmt"
	"net/http"
)

func IndexHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}

func main() {
	http.HandleFunc("/", IndexHandle)
	http.ListenAndServe("127.0.0.1:8000", nil)
}
