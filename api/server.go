package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//Whitelisted url params for filtered api endpoint
var DogTableColumns = []string{
    "id", "name", "origin", "type", "unique_feature", "friendly_rating",
    "life_span", "size", "grooming_needs", "exercise_requirements",
    "good_with_children", "intelligence_rating", "shedding_level",
    "health_issues_risk", "average_weight", "training_difficulty",
}

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
	redisAddr := os.Getenv("REDIS_ADDR")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisPass := os.Getenv("REDIS_PASSWORD")
	dbURL := os.Getenv("DB_URL")

	if redisAddr == "" || err != nil {
		log.Fatal("Redis Client Missing Keys")
	}
    if dbURL == "" {
        log.Fatal("DB_URL environment variable not set")
    }

	redisClient, err := NewRedisClient(redisAddr, redisDB, redisPass)
    if err != nil {
        log.Fatal("Redis client could not be made:", err)
    }

    dbClient, err := NewPostgreSQLClient(dbURL)
    if err != nil {
        log.Fatal("PostgreSQL client could not be made:", err)
    }

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

        exists, err := AddNewToken(r.Context(), redisClient, signedToken)
        if err != nil || !exists {
            http.Error(w, "Failed to add token to redis", http.StatusInternalServerError)
            }

        key_exists, key_err := CheckIfTokenExists(r.Context(), redisClient, signedToken)
        if key_err != nil || !key_exists {
            http.Error(w, "Token not found or error checking token", http.StatusInternalServerError)
            return
        }

        fmt.Println("Key found in Redis")

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(map[string]string{"token": signedToken}); err != nil {
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
            return
        }
	})

	http.HandleFunc("/api/dog-breeds", func(w http.ResponseWriter, r *http.Request){
		enableOpenCors(&w)

        pageStr := r.URL.Query().Get("page")
        limitStr := r.URL.Query().Get("limit")
        page := 1
        limit := 100
        if pageStr != "" {
            fmt.Sscanf(pageStr, "%d", &page)
        }

        if limitStr != "" {
            fmt.Sscanf(limitStr, "%d", &limit)
        }

        authHeader := r.Header.Get("Authorization")

        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            http.Error(w, "Authorization Bearer token missing or improperly formatted", http.StatusInternalServerError)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")

        key_exists, key_err := CheckIfTokenExists(r.Context(), redisClient, token)
        if key_err != nil || !key_exists {
            http.Error(w, "Token not found or error checking token", http.StatusInternalServerError)
            return
        }

        results, db_err := DogBreedsPaginated(dbClient, page, limit)
        if db_err != nil || len(results) < 1 {
            http.Error(w, "Error retrieving data please double check pagination values", http.StatusInternalServerError)
            return
        }

        success, error_message := UpdateTokenUsage(r.Context(), redisClient, token)
        if error_message != "" || !success {
            http.Error(w, error_message, http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(map[string][]Dog{"dogs": results}); err != nil {
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
            return
        }
	})



	http.HandleFunc("/api/dog-breeds/search/{id}", func(w http.ResponseWriter, r *http.Request){
		enableOpenCors(&w)
        idUrl := r.PathValue("id")
        id := 0

        if idUrl != "" {
            fmt.Sscanf(idUrl, "%d", &id)
        }


        authHeader := r.Header.Get("Authorization")

        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            http.Error(w, "Authorization Bearer token missing or improperly formatted", http.StatusInternalServerError)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")

        key_exists, key_err := CheckIfTokenExists(r.Context(), redisClient, token)
        if key_err != nil || !key_exists {
            http.Error(w, "Token not found or error checking token", http.StatusInternalServerError)
            return
        }

        result, db_err := DogById(dbClient, id)
        if db_err != nil {
            http.Error(w, "Error retrieving data please double check id value", http.StatusInternalServerError)
            return
        }

        success, error_message := UpdateTokenUsage(r.Context(), redisClient, token)
        if error_message != "" || !success {
            http.Error(w, error_message, http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(map[string]Dog{"dogs": result}); err != nil {
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
            return
        }
	})

	http.HandleFunc("/api/dog-breeds/filter", func(w http.ResponseWriter, r *http.Request){
		enableOpenCors(&w)

        //Grab all the params most likely in a helper function due to length and amount of params

        authHeader := r.Header.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            http.Error(w, "Authorization Bearer token missing or improperly formatted", http.StatusInternalServerError)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        key_exists, key_err := CheckIfTokenExists(r.Context(), redisClient, token)
        if key_err != nil || !key_exists {
            http.Error(w, "Token not found or error checking token", http.StatusInternalServerError)
            return
        }

        //Db Query for the data using a map ds, make sure to use paramarertized queries

      I  success, error_message := UpdateTokenUsage(r.Context(), redisClient, token)
        if error_message != "" || !success {
            http.Error(w, error_message, http.StatusInternalServerError)
            return
        }

        //Return all data in an array of dog data like the first endpoint

        

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(map[string][]Dog{"dogs": results}); err != nil {
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
            return
        }
	})

	log.Fatal(http.ListenAndServe(PORT, nil))
}

