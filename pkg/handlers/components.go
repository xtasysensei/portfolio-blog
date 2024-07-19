package handlers

import (
	"portfolio-blog/views/templates"

	"github.com/a-h/templ"
)

func Index() *templ.ComponentHandler {
	indexComponent := templates.Index()
	return templ.Handler(indexComponent)
}

func ShowBlogPost() *templ.ComponentHandler {
	showBlogComponent := templates.ShowBlog()
	return templ.Handler(showBlogComponent)
}
func AdminLogin() *templ.ComponentHandler {
	adminComponent := templates.AdminLogin()
	return templ.Handler(adminComponent)
}
func CreatePostForm() *templ.ComponentHandler {
	createPostComponent := templates.CreatePost()
	return templ.Handler(createPostComponent)
}
