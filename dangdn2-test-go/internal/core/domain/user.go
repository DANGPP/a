package domain

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Pass string `json:"pass"`
}

type UserRepository interface {
	GetAll() ([]User, error)
	Create(user User) error
	FindByNameAndPass(name, pass string) (User, error)
}
