package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codeinbit/currency-converter/rates"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	r := mux.NewRouter()

	var err error
	err = godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file : %v", err)
	} else {
		fmt.Println("Success loading .env file")
	}

	rates.InitCache()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/rate/{from}/{to}", rates.UnitRate).Methods("GET")
	api.HandleFunc("/rate/{from}/{to}/{amount}", rates.RateOnAmount).Methods("GET")

	fmt.Printf("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
