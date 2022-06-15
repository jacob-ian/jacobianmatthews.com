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
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Message: "Success"})
	})
	mux.Handle("/", loggerMiddleware(handler))
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
		log.Printf("%v %v %v", r.Method, r.URL.Path, r.UserAgent())
		next.ServeHTTP(w, r)
	})
}
