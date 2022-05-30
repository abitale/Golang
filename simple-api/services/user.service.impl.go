package services

import (
	"context"
	"errors"

	"example.com/simple-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(userCollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{
		userCollection: userCollection,
		ctx:            ctx,
	}
}

func (u *UserServiceImpl) CreateUser(user *models.RegisterUser) error {
	query := bson.D{bson.E{Key: "id", Value: user.ID}}
	checkIfExist := u.userCollection.FindOne(u.ctx, query)
	if checkIfExist.Err() == nil {
		return errors.New("user id already exist")
	}

	query = bson.D{bson.E{Key: "email", Value: user.Email}}
	checkIfExist = u.userCollection.FindOne(u.ctx, query)
	if checkIfExist.Err() == nil {
		return errors.New("email already exist")
	}

	_, err := u.userCollection.InsertOne(u.ctx, user)
	return err
}

func (u *UserServiceImpl) LoginUser(user *models.LoginUser) error {
	var checkUser models.LoginUser
	query := bson.D{bson.E{Key: "email", Value: user.Email}}
	checkIfExist := u.userCollection.FindOne(u.ctx, query).Decode(&checkUser)
	if checkIfExist != nil {
		return errors.New("invalid login request")
	}

	if user.Password != checkUser.Password {
		return errors.New("invalid login request")
	}

	return checkIfExist
}
