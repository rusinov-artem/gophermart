package auth

type Storage interface {
	FindToken(token string) (string, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Auth(token string) (string, error) {
	login, err := s.storage.FindToken(token)
	if err != nil {
		return "", err
	}

	return login, nil
}
