package models

import (
	"github.com/go-openapi/strfmt"
	"time"
)

type PostUpdateRequestInput struct {
	Message string `json:"message"`
}

type PostCreateRequestInput struct {
	Parent uint64 `json:"parent"`
	Author string `json:"author"`
	Message string `json:"message"`
	Created time.Time `json:"created"`
}

type PostDetails struct {
	Post *PostModel `json:"post,omitempty"`
	Author *User `json:"author,omitempty"`
	Thread *ThreadModel `json:"thread,omitempty"`
	Forum *Forum `json:"forum,omitempty"`
}

type PostModel struct {
	ID uint64 `json:"id"`
	Parent uint64 `json:"parent"`
	Author string `json:"author"`
	Message string `json:"message"`
	IsEdited bool `json:"isEdited"`
	Forum string `json:"forum"`
	Thread uint64 `json:"thread"`
	Created strfmt.DateTime `json:"created"`
}

