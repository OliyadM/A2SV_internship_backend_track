package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type Task struct {
	ID          string    `json:"id" bson:"_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	DueDate     time.Time `json:"due_date" bson:"due_date"`
	Status      string    `json:"status" bson:"status"`
}


type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"` 
	Role     string `json:"role" bson:"role"`
}


type TaskRepository interface {
	AddTask(ctx context.Context, task Task) (string, error)
	GetAllTasks(ctx context.Context) ([]Task, error)
	GetTaskByID(ctx context.Context, id string) (*Task, error)
	UpdateTask(ctx context.Context, id string, task Task) error
	DeleteTask(ctx context.Context, id string) error

}


type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	FindUserByUsername(ctx context.Context, username string) (*User, error)
	PromoteUser(ctx context.Context, username string) error
	IsFirstUser(ctx context.Context) (bool, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
	
}


type TaskUsecase interface {
	AddTask(ctx context.Context, task Task) (string, error)
	GetAllTasks(ctx context.Context) ([]Task, error)
	GetTaskByID(ctx context.Context, id string) (*Task, error)
	UpdateTask(ctx context.Context, id string, task Task) error
	DeleteTask(ctx context.Context, id string) error
}


type UserUsecase interface {
	Register(ctx context.Context, user User) error
	Login(ctx context.Context, username, password string) (string, error)
	PromoteUser(ctx context.Context, username string) error
	GetAllUsers(ctx context.Context) ([]*User, error)
}


type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashed, plain string) error
}


type JWTService interface {
	GenerateToken(userID, username, role string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}
