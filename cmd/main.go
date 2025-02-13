package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"forum/cmd/routers"
	"forum/database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		db.Close()
		fmt.Println("\nServer closed succesfully")
		os.Exit(0)
	}()

	routers.Router(db)
	fmt.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
