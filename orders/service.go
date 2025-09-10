package orders

import (
	"fmt"

	"github.com/irreal/order-packs/models"
)

type Service struct {
	MaxOrderItemCount int
}

func NewService(maxOrderItemCount int) *Service {
	return &Service{
		MaxOrderItemCount: maxOrderItemCount,
	}
}

func (s *Service) CalculatePacks(order *models.Order) error {
	if order.RequestedItemCount <= 0 {
		return fmt.Errorf("%w: Item count has to be greater than 0", InvalidOrderItemCountError)
	}
	if order.RequestedItemCount > s.MaxOrderItemCount {
		return fmt.Errorf("%w: Item count has to be less than or equal to %d", InvalidOrderItemCountError, s.MaxOrderItemCount)
	}
	return nil
}
