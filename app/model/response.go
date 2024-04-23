package model

type ErrorResponse struct {
	Error string `json:"error"`
}

type IdResponse struct {
	Id int `json:"id"`
}

type TaskResponse struct {
	Tasks []Task `json:"tasks"`
}
