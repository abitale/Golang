package main

import (
	"context"
	"fmt"
	"log"

	"example.com/simple-api/controllers"
	"example.com/simple-api/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server         *gin.Engine
	mailService    services.MailService
	mailController controllers.MailController
	ctx            context.Context
	mailCollection *mongo.Collection
	mongoClient    *mongo.Client
	err            error
)

func init() {
	ctx = context.TODO()

	mongoConn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err := mongo.Connect(ctx, mongoConn)
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
	mailController = controllers.New(mailService)
	server = gin.Default()
}

func main() {
	defer mongoClient.Disconnect(ctx)

	basepath := server.Group("/v1")
	mailController.MailRoutes(basepath)

	log.Fatal(server.Run(":8080"))
}
