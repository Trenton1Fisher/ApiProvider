package main

import (
	"fmt"
	"log"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main(){
	const PORT = ":5173"

	fmt.Println("Server is running on port:", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}