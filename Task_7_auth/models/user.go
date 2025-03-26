package models

type User struct {
    ID       string `json:"id" bson:"_id"`
    Username string `json:"username" bson:"username" validate:"required,min=3,max=32"`
    Password string `json:"password" bson:"password" validate:"required,min=6"`
    Role     string `json:"role" bson:"role"`
}