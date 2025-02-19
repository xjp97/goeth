package models

type User struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
}
