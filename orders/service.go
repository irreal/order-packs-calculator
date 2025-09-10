package orders

import (
	"fmt"
	"time"

	"github.com/irreal/order-packs/models"
)

type Service struct {
	MaxOrderItemCount int
	repo              OrderRepository
}

type OrderRepository interface {
	SaveOrder(order *models.Order) error
	GetLast10Orders() ([]*models.Order, error)
}

func NewService(maxOrderItemCount int, repo OrderRepository) *Service {
	return &Service{
		MaxOrderItemCount: maxOrderItemCount,
		repo:              repo,
	}
}

func (s *Service) CreateOrder(orderRequest models.OrderRequest, availablePacks []models.Pack) (*models.Order, error) {
	if orderRequest.ItemCount <= 0 {
		return nil, fmt.Errorf("%w: Item count has to be greater than 0", InvalidOrderItemCountError)
	}
	if orderRequest.ItemCount > s.MaxOrderItemCount {
		return nil, fmt.Errorf("%w: Item count has to be less than or equal to %d", InvalidOrderItemCountError, s.MaxOrderItemCount)
	}

	packsCalculation, err := CalculatePack(availablePacks, orderRequest.ItemCount)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", OrderCalculationError, err)
	}

	order := &models.Order{
		RequestedItemCount: orderRequest.ItemCount,
		ShippedItemCount:   packsCalculation.TotalItems,
		Packs:              packsCalculation.Packs,
		Status:             models.OrderStatusNew,
		CreatedAt:          time.Now(),
	}

	// persist order to repo
	if err := s.repo.SaveOrder(order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}

func (s *Service) GetLast10Orders() ([]*models.Order, error) {
	return s.repo.GetLast10Orders()
}
