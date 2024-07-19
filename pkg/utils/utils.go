package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"portfolio-blog/pkg/models"
)

func GetBlogData(w http.ResponseWriter, r *http.Request) {
	var blog models.Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	blogPostID := blog.PostID
	blogTitle := blog.Title
	blogContent := blog.Content
	blogSlug := blog.Slug
	blogTags := blog.Tags

	w.Write([]byte("Blog succcessfully posted"))

	log.Println(blogPostID)
	log.Println(blogTitle)
	log.Println(blogContent)
	log.Println(blogSlug)
	log.Println(blogTags)

}
