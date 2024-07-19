package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"portfolio-blog/passwordhashing"
	"portfolio-blog/pkg/database"
	"portfolio-blog/pkg/models"

	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func AuthAdmin(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into an Admin struct
	var admin models.Admin
	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Changed to BadRequest
		return
	}

	verifyUsername := admin.Username
	verifyPassword := admin.Password

	// Connect to the database
	db := ConnectDB(w)

	// Retrieve stored username and password hash from the database
	storedUsername, storedPasswordHash, err := database.RetrieveUserDB(db, verifyUsername)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Verify the provided password with the stored hash
	match := passwordhashing.VerifyPassword(verifyPassword, storedPasswordHash)

	// Check if credentials match
	if storedUsername != verifyUsername || !match {
		log.Printf("Login failed for username: %s. Username or Password invalid", verifyUsername)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Set the cookie and redirect to the target page
	SetCookie(storedUsername, w)
	redirectTarget := "/create-post"
	http.Redirect(w, r, redirectTarget, http.StatusFound) // Changed status code to StatusFound
	log.Printf("Login successful for username: %s", verifyUsername)
}

func CreateAdmin(w http.ResponseWriter, r *http.Request) {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	passwordhash, err := passwordhashing.HashPassword(password)
	if err != nil {
		http.Error(w, "Password hashing failed: "+err.Error(), http.StatusInternalServerError)
	}

	db := ConnectDB(w)

	if username != "" && passwordhash != "" {
		query := "INSERT INTO admin(username, password_hash) VALUES ($1, $2)"
		if _, err := db.Exec(query, username, passwordhash); err != nil {
			http.Error(w, "Query to database failed: "+err.Error(), http.StatusInternalServerError)
			log.Printf("An error occurred while executing query: %v", err)
			return
		}

		log.Println("Data successfully inserted")
		log.Println("Admin creation successful")

	} else {
		log.Println("Username, email, and password hash must not be empty")
	}

}

// for POST
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	ClearCookie(response)
	http.Redirect(response, request, "/", 302)
}

// Cookie

func SetCookie(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("cookie", value); err == nil {
		cookie := &http.Cookie{
			Name:  "cookie",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func ClearCookie(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "cookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("cookie"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("cookie", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}
