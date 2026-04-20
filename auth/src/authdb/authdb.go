package authdb

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json: "user_id"`
	Name     string `json:"user_name"`
	Password string `json:"user_password"`
}

func Connect(dbRoot string, dbPassword string, dbHost string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/", dbRoot, dbPassword, dbHost))
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}

func CreateDB(db *sql.DB) {
	cmd, err := db.Query("CREATE DATABASE IF NOT EXISTS auth")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer cmd.Close()
}

func CreateTables(db *sql.DB) {
	cmd, err := db.Query("CREATE TABLE IF NOT EXISTS auth.users (user_id int AUTO_INCREMENT, user_name char(50) NOT NULL, user_password char(128), PRIMARY KEY(user_id));")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer cmd.Close()
}

func InsertUser(db *sql.DB, user User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO auth.users (user_name, user_password) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Name, string(hash))
	return err
}

func GetUserByName(user_name string, db *sql.DB) (User, error) {
	var user User
	row := db.QueryRow("SELECT user_id, user_name, user_password FROM auth.users WHERE user_name = ?", user_name)
	err := row.Scan(&user.ID, &user.Name, &user.Password)
	if err == sql.ErrNoRows {
		return User{}, nil
	}
	return user, err
}

func CreateUser(db *sql.DB, u User) (bool, error) {
	user, err := GetUserByName(u.Name, db)
	if err != nil {
		return false, err
	}
	if user != (User{}) {
		return false, nil
	}
	err = InsertUser(db, u)
	if err != nil {
		return false, err
	}
	return true, nil
}
