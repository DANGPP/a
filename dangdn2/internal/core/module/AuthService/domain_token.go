package AuthService

type Token struct {
	UUID      string `gorm:"primaryKey";type:"uuid"`
	Service   string
	Exp       int64
	Iat       int64
	HashToken string
	Status    string
}
