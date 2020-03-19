package main

import (
	"log"
	"net/http"
	"api_stats/router"
)

func main() {

	router := router.Router()

	//fmt.Println("Démarrage du serveur port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}