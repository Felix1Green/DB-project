package models


type PostUpdateRequestInput struct {
	Message string
}

type PostCreateRequestInput struct {
	Parent uint64
	Author string
	Message string
}

type PostModel struct {
	ID uint64
	Parent uint64
	Author string
	Message string
	IsEdited bool
	Forum string
	Thread uint64
	Created string
}

