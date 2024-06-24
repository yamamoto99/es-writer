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
	err := db.QueryRow("SELECT id, username, email, created_at FROM users WHERE username = ?", "testuser").Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
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
	if cognitoRegion == "" || clientId == "" || jwksURL == "" {
		log.Fatalf("congnitまたはgeminiの環境変数が設定されていません")
		fmt.Println(cognitoRegion, clientId, jwksURL, apiKey)
	}

	fmt.Println("db inf: "+os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
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
