package models

import "time"

type PostUpdateRequestInput struct {
	Message string `json:"message"`
}

type PostCreateRequestInput struct {
	Parent uint64 `json:"parent"`
	Author string `json:"author"`
	Message string `json:"message"`
	Created time.Time `json:"created"`
}

type PostModel struct {
	ID uint64 `json:"id"`
	Parent uint64 `json:"parent"`
	Author string `json:"author"`
	Message string `json:"message"`
	IsEdited bool `json:"isEdited"`
	Forum string `json:"forum"`
	Thread uint64 `json:"thread"`
	Created string `json:"created"`
}

