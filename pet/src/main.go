package main

import (
	orm "databases"
	"net/http"
	"routes"
)

func main() {

	defer orm.Eloquent.Close()

	router := router.InitRouter()

	http.ListenAndServe(":8080", router)

}