package main

import (
	"ensemble"
	"ensemble/gui"
	"log"
	"net/http"
	"os"
)

func main() {
	c, err := ensemble.ConfigFromEnv(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	a, err := ensemble.New(log.Default()).NewAgent(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("GUI listening on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", gui.New(a)))
}
