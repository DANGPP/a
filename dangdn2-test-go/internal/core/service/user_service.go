package service

import (
	"dangdn2-test-go/internal/core/domain"
)

type UserService struct {
	Repo domain.UserRepository
}

func (s *UserService) GetAllUsers() ([]domain.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) Register(user domain.User) error {
	return s.Repo.Create(user)
}

func (s *UserService) Login(name, pass string) (domain.User, error) {
	return s.Repo.FindByNameAndPass(name, pass)
}
