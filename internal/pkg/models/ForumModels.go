package models


type Forum struct{
	Title string `json:"title"`
	User string `json:"user"`
	Slug string `json:"slug"`
	Posts uint64 `json:"posts"`
	Threads uint64 `json:"threads"`
}

type ForumRequestInput struct {
	Title string `json:"title"`
	User string `json:"user"`
	Slug string `json:"slug"`
}

