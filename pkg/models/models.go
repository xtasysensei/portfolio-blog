package models

type Blog struct {
	PostID  int      `json:"postid"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Slug    string   `json:"slug"`
	Tags    []string `json:"tags"`
}

type Admin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
