package service

import (
	"github.com/rakaascode/server-antrian-go.git/models"
	"github.com/rakaascode/server-antrian-go.git/repository"
)

type UserService interface {
	GetAll() ([]models.User, error)
	GetByID(id uint) (models.User, error)
	Create(user models.User) (models.User, error)
	Update(id uint, user models.User) (models.User, error)
	Delete(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{r}
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetByID(id uint) (models.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) Create(user models.User) (models.User, error) {
	return s.repo.Create(user)
}

func (s *userService) Update(id uint, user models.User) (models.User, error) {
	user.ID = id
	return s.repo.Update(user)
}

func (s *userService) Delete(id uint) error {
	return s.repo.Delete(id)
}