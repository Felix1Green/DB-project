package models

type ThreadRequestInput struct {
	Title string
	Author string
	Message string
	Created string
}

type ThreadModel struct {
	ID uint64
	Title string
	Author string
	Message string
	Created string
	Votes uint64
	Forum string
}

type ThreadUpdateInput struct {
	Title string
	Message string
}

type ThreadVoteInput struct {
	Nickname string
	Voice int
}