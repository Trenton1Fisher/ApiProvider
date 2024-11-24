package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main(){
	const PORT = ":9010"
	redisAddr := os.Getenv("REDIS_ADDR")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisPass := os.Getenv("REDIS_PASSWORD")
	dbURL := os.Getenv("DB_URL")

	if redisAddr == "" || err != nil {
		log.Fatal("Redis Client Missing Keys;")
	}
    if dbURL == "" {
        log.Fatal("DB_URL environment variable not set")
    }

	redisClient, err := NewRedisClient(redisAddr, redisDB, redisPass)
    if err != nil {
        log.Fatal("Redis client could not be made:", err)
    }

    // Initialize PostgreSQL client
    dbClient, err := NewPostgreSQLClient(dbURL)
    if err != nil {
        log.Fatal("PostgreSQL client could not be made:", err)
    }


	http.HandleFunc("/api/get-token", func(w http.ResponseWriter, r *http.Request){
		enableCors(&w)
	})

	http.HandleFunc("/api/dog-breeds", func(w http.ResponseWriter, r *http.Request){
		enableCors(&w)
	})

	http.HandleFunc("/api/dog-breeds/search/id", func(w http.ResponseWriter, r *http.Request){
		enableCors(&w)
	})

	http.HandleFunc("/api/dog-breeds/filter", func(w http.ResponseWriter, r *http.Request){
		enableCors(&w)
	})

	//log.Fatal(http.ListenAndServe(PORT, nil))
}