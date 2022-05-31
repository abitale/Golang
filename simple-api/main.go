package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"example.com/simple-api/auth"
	"example.com/simple-api/controllers"
	"example.com/simple-api/services"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	mailCollection *mongo.Collection
	mailService    services.MailService
	mailController controllers.MailController
	userCollection *mongo.Collection
	userService    services.UserService
	userController controllers.UserController
	tokenAuth      *jwtauth.JWTAuth
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
	userController = controllers.NewUser(userService)
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

func Authen(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if auth.ValidateToken(token) != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	defer mongoClient.Disconnect(ctx)

	server := chi.NewRouter()
	server.Use(middleware.Logger)
	server.Use(middleware.Recoverer)
	server.Route("/users", func(r chi.Router) {
		r.Post("/register", userController.CreateUser)
		r.Post("/login", userController.LoginUser)
	})
	server.Route("/mails", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(Authen)
		r.Post("/", mailController.CreateMail)
		r.Get("/", mailController.GetAll)
		r.Get("/{id}", mailController.GetMail)
		r.Put("/{id}", mailController.UpdateMail)
		r.Delete("/{id}", mailController.DeleteMail)
	})
	http.ListenAndServe(":8080", server)
}
