package repository

import "lld/lld-2/bookmyshow/models"

type UserRepo struct {
	UserList map[int]models.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		UserList: make(map[int]models.User),
	}
}

func (ur *UserRepo) SaveUser(u models.User) error {
	ur.UserList[u.Id] = u
	return nil
}

func (ur *UserRepo) FindUserById(id int) models.User {
	return ur.UserList[id]
}
