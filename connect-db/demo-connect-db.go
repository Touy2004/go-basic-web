package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func connectDb() (db *sql.DB) {
	var (
		userDB string = "root"
		passDB string = ""
		hostDB string = "localhost"
		portDB string = "3306"
		dbName string = "coursesdb"
	)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userDB, passDB, hostDB, portDB, dbName))

	if err != nil {
		fmt.Println("failed to connect database")
	} else {
		fmt.Println("connected to database")
	}
	return db
}

func query(db *sql.DB) {
	var inputId int
	fmt.Println("Which id do you want to see?")
	fmt.Print("id: ")
	fmt.Scan(&inputId)
	var (
		id         int
		coursename string
		price      float64
		instructor string
	)

	query := "SELECT id, coursename, price, instructor FROM onlinecourse WHERE id = ?"
	if err := db.QueryRow(query, inputId).Scan(&id, &coursename, &price, &instructor); err != nil {
		log.Fatal(err)
	}
	fmt.Println(id, coursename, price, instructor)
}

func creatingTable(db *sql.DB) {
	query := `CREATE TABLE users (id INT AUTO_INCREMENT, username TEXT NOT NULL, password TEXT NOT NULL, creata_at DATETIME, PRIMARY KEY (id))`
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func InsertUser(db *sql.DB){
	var username string
	var password string 
	fmt.Scan(&username)
	fmt.Scan(&password)
	createAt := time.Now()

	query := `INSERT INTO users (username, password, creata_at) VALUES (?, ?, ?)`

	result, err := db.Exec(query, username, password, createAt)

	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	fmt.Println(id)
}

func DeleteUser(db *sql.DB) {
	var id int
	fmt.Scan(&id)
	query := `DELETE FROM users WHERE id = ?`
	if _, err := db.Exec(query, id); err != nil {
		log.Fatal(err)
	}
}

func main() {
	db := connectDb()
	InsertUser(db)
	DeleteUser(db)
}
