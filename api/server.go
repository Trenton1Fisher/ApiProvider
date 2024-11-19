package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main(){
	const PORT = ":9010"

	dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatal("DB_URL environment variable not set")
    }


	fmt.Println(dbURL)
	//log.Fatal(http.ListenAndServe(PORT, nil))
}