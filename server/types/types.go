package types

type InsertUserParams struct {
	Password string
	Username string
	Email    string
}

type InsertTodoParams struct {
	Text   string
	Done   bool
	ListId int
	UserId int
}

type UpdateTodoParams struct {
	ID     *int    `db:"id"`
	Done   *bool   `db:"done"`
	Text   *string `db:"text"`
	UserId *int
}

type GetTodosParams struct {
	TodoId *int `db:"id"`
	ListId *int `db:"todo_list_id"`
	UserId *int
}

type Todo struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Done   bool   `json:"done"`
	ListId int    `json:"todo_list_id"`
}
