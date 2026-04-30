package cabang

type CabangService interface {
	GetAll() ([]Cabang, error)
	GetByID(id uint) (Cabang, error)
	Create(req CreateCabangRequest) (Cabang, error)
	Update(id uint, req UpdateCabangRequest) (Cabang, error)
	Delete(id uint) error
}

type cabangService struct {
	repo CabangRepository
}

func NewCabangService(r CabangRepository) CabangService {
	return &cabangService{r}
}

func (s *cabangService) GetAll() ([]Cabang, error) {
	return s.repo.FindAll()
}

func (s *cabangService) GetByID(id uint) (Cabang, error) {
	return s.repo.FindByID(id)
}

func (s *cabangService) Create(req CreateCabangRequest) (Cabang, error) {
	return s.repo.Create(Cabang{
		Nama:      req.Nama,
		Alamat:    req.Alamat,
		Kota:      req.Kota,
		NoTelp:    req.NoTelp,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
}

func (s *cabangService) Update(id uint, req UpdateCabangRequest) (Cabang, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return Cabang{}, err
	}
	if req.Nama != "" {
		c.Nama = req.Nama
	}
	if req.Alamat != "" {
		c.Alamat = req.Alamat
	}
	if req.Kota != "" {
		c.Kota = req.Kota
	}
	if req.NoTelp != "" {
		c.NoTelp = req.NoTelp
	}
	if req.Latitude != 0 {
		c.Latitude = req.Latitude
	}
	if req.Longitude != 0 {
		c.Longitude = req.Longitude
	}
	return s.repo.Update(c)
}

func (s *cabangService) Delete(id uint) error {
	return s.repo.Delete(id)
}
