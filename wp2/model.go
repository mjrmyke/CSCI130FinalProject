package finalproject

type User struct {
	Email    string //email is the prim key, has to be unique
	Password string
}

type Session struct {
	User    User
	id      string
	logged   bool
	Message string
}
