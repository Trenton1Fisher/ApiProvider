package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "https://trenton1fisher.github.io/ApiProvider")
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

		claims := &Claims {
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer: "Public-Dog-Api",
				IssuedAt: jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		if err != nil {
			http.Error(w, "Error signing token", http.StatusInternalServerError)
			return
		}

		fmt.Println(token)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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