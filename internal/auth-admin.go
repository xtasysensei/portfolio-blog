package internal

import (
	"log"
	"net/http"
	"os"
	"portfolio-blog/passwordhashing"
	"portfolio-blog/pkg/database"

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
	// Read form values
	verifyUsername := r.FormValue("username")
	verifyPassword := r.FormValue("password")

	// Connect to the database
	db := ConnectDB(w)

	// Retrieve stored username and password hash from the database
	storedUsername, storedPasswordHash, err := database.RetrieveUserDB(db, verifyUsername)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	// Verify the provided password with the stored hash
	match := passwordhashing.VerifyPassword(verifyPassword, storedPasswordHash)

	// Check if credentials match
	if storedUsername != verifyUsername || !match {
		log.Printf("Login failed for username: %s. Username or Password invalid", verifyUsername)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized) // Changed to Unauthorized
		return
	}

	// Set the cookie
	SetCookie(storedUsername, w)

	// Log successful login
	log.Printf("Login successful for username: %s", verifyUsername)

	// Redirect to the target page
	// data := map[string]interface{}{} // Replace with actual data if needed
	// err = templ.Render(w, "create-post", data)
	// if err != nil {
	// 	log.Printf("Error rendering template: %v", err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }
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
func SetCookie(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("cookie", value); err == nil {
		cookie := &http.Cookie{
			Name:     "cookie",
			Value:    encoded,
			Path:     "/",
			HttpOnly: true, // Prevent client-side access
			Secure:   true, // Set to true if using HTTPS
		}
		http.SetCookie(response, cookie)
	} else {
		log.Printf("Error encoding cookie: %v", err)
	}
}

// ClearCookie removes the authentication cookie
func ClearCookie(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "cookie",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,   // Remove the cookie
		HttpOnly: true, // Prevent client-side access
		Secure:   true, // Set to true if using HTTPS
	}
	http.SetCookie(response, cookie)
}

// GetUserName retrieves the username from the cookie
func GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("cookie"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("cookie", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		} else {
			log.Printf("Error decoding cookie: %v", err)
		}
	}
	return userName
}

// LogoutHandler handles user logout by clearing the cookie and redirecting
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	ClearCookie(response)
	log.Println("User logged out")
	http.Redirect(response, request, "/", http.StatusSeeOther) // Use StatusSeeOther (303) for redirection after POST
}
