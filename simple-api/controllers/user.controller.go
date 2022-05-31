package controllers

import (
	"encoding/json"
	"net/http"

	"example.com/simple-api/auth"
	"example.com/simple-api/models"
	"example.com/simple-api/services"
)

type UserController struct {
	UserService services.UserService
}

func NewUser(userService services.UserService) UserController {
	return UserController{
		UserService: userService,
	}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := uc.UserService.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Write([]byte("success"))
	w.WriteHeader(201)
}

func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := uc.UserService.LoginUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokenString, err := auth.GenerateJWT(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("token: " + tokenString))
	w.WriteHeader(201)
}
