package packs

import "github.com/irreal/order-packs/models"

type Service struct {
	repo PackRepository
}

type PackRepository interface {
	GetPacks() (models.Packs, error)
	SavePacks(packs models.Packs) error
}

func NewService(repo PackRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetPacks() (models.Packs, error) {
	return s.repo.GetPacks()
}

func (s *Service) SavePacks(packs models.Packs) error {
	return s.repo.SavePacks(packs)
}
