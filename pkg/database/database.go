package database

import (
	"database/sql"
	"fmt"
)

func RetrieveUserDB(db *sql.DB, queryUsername string) (string, string, error) {
	var getUsername string
	var getPassword string

	// Corrected SQL query with a comma between columns
	query := "SELECT username, password_hash FROM admin WHERE username = $1"

	err := db.QueryRow(query, queryUsername).Scan(&getUsername, &getPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("user not found")
		}
		return "", "", fmt.Errorf("error retrieving user data: %v", err)
	}

	return getUsername, getPassword, nil
}
