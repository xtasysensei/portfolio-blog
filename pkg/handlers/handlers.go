package handlers

import (
	"net/http"
	"portfolio-blog/pkg/utils"
)

func CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	utils.GetBlogData(w, r)
}
