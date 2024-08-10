package service

import "github.com/datnguyen210/go-muji/internal/repository"

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepository: repository.NewUserRepository(),
	}
}

func (us UserService) ReadUserInfo() string {
	return us.userRepository.ReadUserInfo()
}
