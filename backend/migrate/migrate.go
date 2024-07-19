package main

import (
	"es-app/db"
	"es-app/model"
	"fmt"
)

func main() {
	dbConn := db.NewDB()
	defer fmt.Println("🟢 Successfully migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{})
}
