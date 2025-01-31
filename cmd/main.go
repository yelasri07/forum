package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/cmd/routers"
	"forum/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	routers.Router(db)
	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
