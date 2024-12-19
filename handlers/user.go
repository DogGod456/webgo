package handlers

type Note struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Color   string `json:"color"`
}

type User struct {
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashedPassword"`
	Notes          []Note `json:"notes"`
}
