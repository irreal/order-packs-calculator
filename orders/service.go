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
	}

	return order, nil
}
