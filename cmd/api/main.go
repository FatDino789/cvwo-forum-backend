package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const port = 8080

type application struct {
	Domain string
	JwtSecret string
}

func main() {
	// set application config
	app := &application{
		Domain: "example.com", 
		JwtSecret: os.Getenv("JWT_SECRET"),
	}

	log.Println("Starting application on port ", port)

	// start a web server
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}