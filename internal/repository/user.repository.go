package repository

import logger "github.com/datnguyen210/go-muji/pkg/util"

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) ReadUserInfo() string {
	log := logger.CreateLogger()
	defer log.Sync()
	log.Info("Called to ReadUserInfo successfully, response returned: ")
	return "Chris"
}
