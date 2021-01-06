package models


type UserRequestBody struct{
	FullName string
	About string
	Email string
}


type User struct{
	Nickname string
	FullName string
	About string
	Email string
}