package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
}

func enableRestrictedCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "https://trenton1fisher.github.io")
}

func enableOpenCors(w * http.ResponseWriter){
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main(){
	const PORT = ":9010"
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
/*	redisAddr := os.Getenv("REDIS_ADDR")
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

	err = dbClient.Ping()
	if err != nil {
		log.Fatal("PostgreSQL ping failed:", err)
	}
	log.Println("Successfully connected to PostgreSQL")
*/

	http.HandleFunc("/api/get-token", func(w http.ResponseWriter, r *http.Request){
		enableRestrictedCors(&w)

		//Don't really feel like doing a crazy system so ill just use times
		claims := &Claims {
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer: "Public-Dog-Api",
				IssuedAt: jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        signedToken, err := token.SignedString(jwtSecret)
        if err != nil {
            http.Error(w, "Error signing token", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(map[string]string{"token": signedToken}); err != nil {
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
            return
        }
	})

	http.HandleFunc("/api/dog-breeds", func(w http.ResponseWriter, r *http.Request){
		enableOpenCors(&w)
	})

	http.HandleFunc("/api/dog-breeds/search/id", func(w http.ResponseWriter, r *http.Request){
		enableOpenCors(&w)
	})

	http.HandleFunc("/api/dog-breeds/filter", func(w http.ResponseWriter, r *http.Request){
		enableOpenCors(&w)
	})

	log.Fatal(http.ListenAndServe(PORT, nil))
}