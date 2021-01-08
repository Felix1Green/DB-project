package models


type UserRequestBody struct{
	FullName string `json:"fullname"`
	About string `json:"about"`
	Email string `json:"email"`
}


type User struct{
	Nickname string `json:"nickname"`
	FullName string `json:"fullname"`
	About string `json:"about"`
	Email string `json:"email"`
}