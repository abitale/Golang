# Golang simple-api
CRUD API with MongoDB and Gin Web Framework
## Requirements:
    1) Installed Latest GO / LTS on your computer
        You can check your GO version by running this in terminal: `go version`
    2) MongoDB running on your computer
        You can check if you have MongoDB running with this command: `mongo`

You can run the program with this command in your terminal (make sure you are in the right directory where the `main.go` located):
### `go run main.go`

## How To:
The API runs on port 8080, so it will be http://localhost:8080/

Here is the available routes:
1) http://localhost:8080/v1/mails/ (Need Authorization on Header)
    - POST = Create New Mail
    - GET = Get All Mail
2) http://localhost:8080/v1/mails/{id} (Need Authorization on Header)
    - PUT = Update Mail by ID
    - GET = Get Mail by ID
    - DELETE = Delete Mail by ID
3) http://localhost:8080/auth/users/register
    - POST = Create New User
4) http://localhost:8080/auth/users/login
    - POST = Login User to get Token
