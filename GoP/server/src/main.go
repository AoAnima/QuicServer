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

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Ty struct{}

func ListenAndServeTLS() {
	err := http.ListenAndServeTLS(":443", "cert/cert.pem", "cert/key.pem", http.HandlerFunc(обработчикЗапроса))
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

func обработчикЗапроса(w http.ResponseWriter, req *http.Request) {

	АнализЗапроса(w, req)
	fmt.Printf(" %s  %s \n", w, *req)
}
