package database

import (
	"database/sql"
	"fmt"
)

func RetrieveUserDB(db *sql.DB, query_username string) (string, string, error) {
	var getusername string
	var getpassword string

	query := "SELECT username password_hash FROM admin WHERE username = $1"
	err := db.QueryRow(query, query_username).Scan(&getusername, &getpassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("user not found")
		}
		return "", "", fmt.Errorf("error retrieving user data: %v", err)
	}

	return getusername, getpassword, nil
}
