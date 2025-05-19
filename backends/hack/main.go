package main

import (
	"fmt"
	"log"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	dump(r)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte("{\"hello\":\"world\"}"))
	if err != nil {
		log.Fatal(err)
	}
}

func dump(r *http.Request) {
	fmt.Printf("%v\n", r)
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8081", nil)
}
