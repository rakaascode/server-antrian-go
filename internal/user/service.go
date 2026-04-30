package user

type UserService interface {
	GetAll() ([]User, error)
	GetByID(id uint) (User, error)
	Create(user User) (User, error)
	Update(id uint, user User) (User, error)
	Delete(id uint) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(r UserRepository) UserService {
	return &userService{r}
}

func (s *userService) GetAll() ([]User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetByID(id uint) (User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) Create(u User) (User, error) {
	return s.repo.Create(u)
}

func (s *userService) Update(id uint, u User) (User, error) {
	u.ID = id
	return s.repo.Update(u)
}

func (s *userService) Delete(id uint) error {
	return s.repo.Delete(id)
}
