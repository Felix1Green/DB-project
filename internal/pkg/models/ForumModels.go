package models


type Forum struct{
	Title string
	User string
	Slug string
	Posts uint64
	Threads uint64
}

type ForumRequestInput struct {
	Title string
	User string
	Slug string
}

