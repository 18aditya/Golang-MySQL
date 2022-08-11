package model

type Posts struct {
	Id          int    `form:"idposts" json:"idposts"`
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
}

type Users struct {
	Id         int    `form:"idusers" json:"idusers"`
	First_name string `form:"First_name" json:"First_name"`
	Last_name  string `form:"Last_name" json:"Last_name"`
	Email      string `form:"Email" json:"Email"`
	CreatedAt  string `form:"CreatedAt" json:"CreatedAt"`
	Posts      []Posts
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Posts
}
