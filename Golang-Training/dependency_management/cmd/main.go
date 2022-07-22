package main

import (
	"github.com/hiteshpattanayak-tw/golangtraining/dependency_management/models"
	"github.com/hiteshpattanayak-tw/golangtraining/dependency_management/repository"
)

func main() {
	userRepo := repository.NewUserRepository()

	user1 := models.GetNewUser("Hitesh", "hitesh@tw.com", 12345)
	user2 := models.GetNewUser("Goutam", "goutam@tw.com", 67890)

	userRepo.AddUser(user1)
	userRepo.AddUser(user2)

	userRepo.PrintAllUserDetails()
}
