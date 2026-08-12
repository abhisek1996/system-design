package services

import (
	"lld/lld-2/bookmyshow/interfaces"
	"lld/lld-2/bookmyshow/models"
)

type UserService struct {
	UserRepo interfaces.UserRepo
}

func (us *UserService) RegisterUser(name string, email string) error {

	user := models.User{
		Name:  name,
		Email: email,
	}
	us.UserRepo.SaveUser(user)
	return nil
}

func (us *UserService) GetUser(id int) (models.User, error) {
	return us.UserRepo.FindUserById(id), nil
}
