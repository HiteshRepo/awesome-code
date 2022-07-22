package repository

import (
	"fmt"
	"github.com/hiteshpattanayak-tw/golangtraining/dependency_management/models"
)

type UserRepository struct {
	users []*models.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{users: make([]*models.User, 0)}
}

func (ur *UserRepository) AddUser(user *models.User) {
	ur.users = append(ur.users, user)
}

func (ur *UserRepository) PrintAllUserDetails() {
	for _,u := range ur.users {
		fmt.Printf("Id: %d, Name: %s, Email: %s\n", u.GetId(), u.GetName(), u.GetEmail())
	}
}
