package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Uuid     string
	Username string
	Password string
	Fname    string
	Lname    string
	Email    string
	Errors map[string]string
}

func saveData(u *User) error {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists users (firstname text, lastname text, username text, email text, password text)")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users (firstname, lastname, username, email, password) values (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(u.Fname, u.Lname, u.Username, u.Email, u.Password)
	tx.Commit()
	return err
}

func userExists(u *User) bool {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()
	var ps, us string
	q, err := db.Query("select username, password from users where username = '" + u.Username + "'")
	if err != nil {
		return false
	}
	for q.Next(){
		q.Scan(&us, &ps)
	}
	pw := bcrypt.CompareHashAndPassword([]byte(ps), []byte(u.Password))
	if us == u.Username && pw == nil {
		return true
	}
	return false
}

func encryptPass(pass string) string {
	p := []byte(pass)
	hashpw, _ := bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
	return string(hashpw)
}