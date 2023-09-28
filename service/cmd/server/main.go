package main

import (
	"fmt"
	"log"

	"net/http"
)

func main() {

	go ListenAndServeTLS()
	log.Printf(" %s", "запустили")
	ListenAndServe()

}

func ListenAndServeTLS() {
	err := http.ListenAndServeTLS(":443", "cert/cert.pem", "cert/key.pem", http.HandlerFunc(requestHandler))
	if err != nil {
		fmt.Println(err)
	}

}
func ListenAndServe() {
	err := http.ListenAndServe(":80", http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "https://localhost:443"+req.RequestURI, http.StatusMovedPermanently)
		}))
	if err != nil {
		fmt.Println(err)
	}
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf(" %s  %s \n", w, *req)
}
