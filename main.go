package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	router := httprouter.New()
	router.GET("/*name", handler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello World, %s!", ps.ByName("name"))
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
