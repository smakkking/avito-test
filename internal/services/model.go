package services

type Service struct {
	bannerStorage Storage
}

type Storage interface {
}

func NewService(storage Storage) *Service {
	return &Service{
		bannerStorage: storage,
	}
}
