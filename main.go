package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"webgo/handlers"
)

//var db *sql.DB

func handleRequest() {
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/register", handlers.RegUser)
	http.HandleFunc("/homeUser", handlers.HomeUser)
	http.HandleFunc("/login", handlers.LoginUser)
	http.HandleFunc("/adminPage", handlers.AdminPage)
	http.HandleFunc("/logout", handlers.Logout)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	var err error
	connStr := "user=postgres password=abobavgo dbname=webgodb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	handlers.SetDB(db)
	handleRequest()
}