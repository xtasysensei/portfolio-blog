package handlers

import (
	"log"
	"net/http"
	"portfolio-blog/internal"
	"strings"
)

func CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	// Extract form values
	blogTitle := r.FormValue("title")
	blogContent := r.FormValue("content")
	blogSlug := r.FormValue("slug")
	blogTags := r.FormValue("tags")

	// Convert tags to slice
	tagSlice := strings.Split(blogTags, ",")

	// Connect to the database
	db := internal.ConnectDB(w)
	defer db.Close()

	// Insert data into the database
	query := `
		 INSERT INTO blog_posts (title, content, slug, tags)
		 VALUES ($1, $2, $3, $4)
	 `
	_, err := db.Exec(query, blogTitle, blogContent, blogSlug, tagSlice)
	if err != nil {
		http.Error(w, "Error inserting data into database:"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.Write([]byte("Blog successfully posted"))
	log.Println(blogTitle)
	log.Println(blogContent)
	log.Println(blogSlug)
	log.Println(tagSlice)

}
