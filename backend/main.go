package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type User struct {
	ID   int
	Name string
}

var db *sql.DB

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler called...")

	var user User
	err := db.QueryRow("SELECT id, name FROM users WHERE name = $1", "Bob").Scan(&user.ID, &user.Name)
	if err != nil {
		fmt.Printf("error in query: %s", err)
		log.Fatal(err)
	}

	fmt.Fprintf(w, "ID: %d, Name: %s\n", user.ID, user.Name)
}

func main() {
	fmt.Println("server started...")
	var err error
	db, err = sql.Open("postgres", "host=db user=postgres password=postgres dbname=testdb sslmode=disable")
	if err != nil {
		fmt.Println("error in db connection")
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
