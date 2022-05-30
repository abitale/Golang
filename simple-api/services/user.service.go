package services

import "example.com/simple-api/models"

type UserService interface {
	CreateUser(*models.RegisterUser) error
	LoginUser(*models.LoginUser) error
}
