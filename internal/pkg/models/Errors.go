package models

import "errors"

var(
	InternalDBError = errors.New("internal service error")
	IncorrectInputParams = errors.New("incorrect input params")
	NoSuchUser = errors.New("user does not exist")
	UserAlreadyExists = errors.New("user already exists")
	ThreadAbsentsError = errors.New("thread absents")
	UserConflict 	   = errors.New("user conflicts")
	ParentPostDoesntExists = errors.New("parent post doesnt exists")
	ForumAlreadyExists	   = errors.New("forum already exists")
	ForumDoesntExists	   = errors.New("forum doesnt exists")
	ThreadUniqueErr		   = errors.New("thread already exists")
	PostDoesntExists 	   = errors.New("post doesnt exists")
	IncorrectPath		   = errors.New("no such file or dir")
	EmptyBody			   = errors.New("empty body")
	ThreadDoesntExist	   = errors.New("thread doesnt exists")
)


var(
	ErrorsStatusCodes = map[error]int {
		InternalDBError: 500,
		IncorrectInputParams: 400,
		NoSuchUser: 404,
		ThreadAbsentsError: 404,
		ParentPostDoesntExists: 409,
		ForumDoesntExists: 404,
		UserAlreadyExists: 409,
		UserConflict: 409,
		EmptyBody: 201,
		ThreadUniqueErr: 409,
	}
)

type ErrorMessage struct {
	Message string `json:"message"`
}