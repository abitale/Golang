package main

import (
	"context"
	"fmt"
	"log"

	"example.com/simple-api/auth"
	"example.com/simple-api/controllers"
	"example.com/simple-api/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server         *gin.Engine
	mailCollection *mongo.Collection
	mailService    services.MailService
	mailController controllers.MailController
	userCollection *mongo.Collection
	userService    services.UserService
	UserController controllers.UserController
	ctx            context.Context
	mongoClient    *mongo.Client
	err            error
)

func init() {
	ctx = context.TODO()

	mongoConn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err = mongo.Connect(ctx, mongoConn)
	if err != nil {
		log.Fatal(err)
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err) //mengecek apakah bisa connect ke database utama
	}

	fmt.Println("mongo connection established")

	mailCollection = mongoClient.Database("maildb").Collection("mails")
	mailService = services.NewMailService(mailCollection, ctx)
	mailController = controllers.NewMail(mailService)

	userCollection = mongoClient.Database("maildb").Collection("users")
	userService = services.NewUserService(userCollection, ctx)
	UserController = controllers.NewUser(userService)

	server = gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "access denied, not authorized"})
			c.Abort()
			return
		}
		err := auth.ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	defer mongoClient.Disconnect(ctx)

	basepath := server.Group("/v1", AuthRequired())
	mailController.MailRoutes(basepath)
	UserController.UserRoutes(server.Group("/auth"))

	log.Fatal(server.Run(":8080"))
}
