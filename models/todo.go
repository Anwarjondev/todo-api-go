package models

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	UserId    int    `json:"user_id"`
}

type TodoModel struct {
	Title string `json:"title"`
}
type UpdateTodoModel struct {
	Title string `json:"title"`
	Completed bool `json:"completed"`
}