package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("hello")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v", r)
		fmt.Fprintf(w, "Hi, from server2")
	})

	log.Printf("registering")
	resp, err := http.Post("http://localhost:9000/register?url="+"http://localhost:8082", "", nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)

	log.Fatal(http.ListenAndServe(":8082", nil))
}
