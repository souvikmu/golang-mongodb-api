package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/souvikmukherjee/mongoapi/router"
)

func main() {
	fmt.Println("mongodb API")
	fmt.Println("server is starting...")
	r := router.Router()

	log.Fatal(http.ListenAndServe(":8000", r))
}
