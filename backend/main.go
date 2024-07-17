package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt string
}

var db *sql.DB

var cognitoRegion string
var clientId string
var jwksURL string
var apiKey string

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler called...")

	var user User
	err := db.QueryRow("SELECT id, username, email, created_at FROM users WHERE username = $1", "testuser").Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		fmt.Printf("error in query: %s", err)
		return
	}

	fmt.Fprintf(w, "ID: %s, Name: %s, Email: %s, CreatedAt: %s\n", user.ID, user.Name, user.Email, user.CreatedAt)
}

func main() {
	fmt.Println("server started...")
	var err error

	cognitoRegion = os.Getenv("COGNITO_REGION")
	clientId = os.Getenv("COGNITO_CLIENT_ID")
	jwksURL = os.Getenv("TOKEN_KEY_URL")
	apiKey = os.Getenv("GOOGLE_API_KEY")
	if cognitoRegion == "" || clientId == "" || jwksURL == "" || apiKey == "" {
		log.Fatalf("congnitまたはgeminiの環境変数が設定されていません")
	}

	var dbHost string = os.Getenv("DB_HOST")
	var dbUser string = os.Getenv("DB_USER")
	var dbPassword string = os.Getenv("DB_PASSWORD")
	var dbName string = os.Getenv("DB_NAME")
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatalf("環境変数が設定されていません:1")
	}

	db, err = sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName))
	if err != nil {
		fmt.Println("error in db connection")
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", handler)
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/checkEmail", checkEmail)
	http.HandleFunc("/resendEmail", resendEmail)
	http.HandleFunc("/welcome", welcome)
	http.HandleFunc("/saveprofile", saveProfile)
	http.HandleFunc("/getAnswers", processQuestionsWithAI)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
