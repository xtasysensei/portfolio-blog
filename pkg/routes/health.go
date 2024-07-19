package routes

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server is up and running!"))
}
