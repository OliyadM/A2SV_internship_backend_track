package infrastructure

import (
	"task_manager/domain"

	"golang.org/x/crypto/bcrypt"
)


type PasswordServiceImpl struct{}

func NewPasswordService() domain.PasswordService {
	return &PasswordServiceImpl{}
}

func (s *PasswordServiceImpl) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func (s *PasswordServiceImpl) ComparePassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
