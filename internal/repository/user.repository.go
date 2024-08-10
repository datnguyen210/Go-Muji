package repository

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) ReadUserInfo() string {
	return "Chris"
}
