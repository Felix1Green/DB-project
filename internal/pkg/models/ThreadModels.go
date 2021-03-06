package models

import "github.com/go-openapi/strfmt"

type ThreadRequestInput struct {
	Title string `json:"title"`
	Author string `json:"author"`
	Message string `json:"message"`
	Created string `json:"created"`
	Slug string `json:"slug"`
}

type ThreadModel struct {
	ID uint64 `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Message string `json:"message"`
	Created strfmt.DateTime `json:"created"`
	Votes int64 `json:"votes"`
	Forum string `json:"forum"`
	Slug string `json:"slug"`
}

type ThreadUpdateInput struct {
	Title string `json:"title"`
	Message string `json:"message"`
}

type ThreadVoteInput struct {
	Nickname string `json:"nickname"`
	Voice int `json:"voice"`
}