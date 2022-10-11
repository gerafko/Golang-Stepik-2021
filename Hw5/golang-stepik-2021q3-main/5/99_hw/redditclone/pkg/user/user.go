package user

type User struct {
	ID       string `json:"id"`
	Login    string `json:"username"`
	Password string `json:"password"`
}

type UserRepo struct {
	data map[string]*User
}
