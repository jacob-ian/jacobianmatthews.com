package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	Message string `json:"message"`
}

var Port int

func main() {
	Port := getPort()
	mux := http.NewServeMux()
	handlerA := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Message: "Success"})
	})
	handlerB := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Yo"))
	})
	mux.Handle("/", loggerMiddleware(handlerB))

	mux.Handle("/api", loggerMiddleware(handlerA))
	log.Printf("Listening on Port %v\n", Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Port), mux))
}

func getPort() int {
	portEnv, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatalln("Missing PORT environment variable")
	}
	port, err := strconv.ParseUint(portEnv, 10, 16)
	if err != nil {
		log.Fatalln("Invalid PORT environment variable")
	}
	return int(port)
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("%v %v %v %v", r.Method, r.URL.Path, w.Header().Values("*"), r.UserAgent())
	})
}
