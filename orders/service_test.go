package orders

import (
	"errors"
	"testing"

	"github.com/irreal/order-packs/models"
)

func TestOrderService_CreateOrder_HappyPath(t *testing.T) {
	tests := []struct {
		name         string
		maxCount     int
		orderRequest models.OrderRequest
		packs        []models.Pack
	}{
		{
			name:         "small order",
			maxCount:     1000000000,
			orderRequest: models.OrderRequest{ItemCount: 1},
			packs:        []models.Pack{250, 500, 1000, 2000, 5000},
		},
		{
			name:         "medium order",
			maxCount:     1000000000,
			orderRequest: models.OrderRequest{ItemCount: 251},
			packs:        []models.Pack{250, 500, 1000, 2000, 5000},
		},
		{
			name:         "large order from main example",
			maxCount:     1000000000,
			orderRequest: models.OrderRequest{ItemCount: 12001},
			packs:        []models.Pack{250, 500, 1000, 2000, 5000},
		},
		{
			name:         "order at max limit",
			maxCount:     1000,
			orderRequest: models.OrderRequest{ItemCount: 1000},
			packs:        []models.Pack{250, 500, 1000, 2000, 5000},
		},
		{
			name:         "order with single pack type",
			maxCount:     1000000000,
			orderRequest: models.OrderRequest{ItemCount: 5},
			packs:        []models.Pack{5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.maxCount)
			order, err := service.CreateOrder(tt.orderRequest, tt.packs)

			// no errors
			if err != nil {
				t.Fatalf("CreateOrder() unexpected error = %v", err)
			}

			// has order
			if order == nil {
				t.Fatal("CreateOrder() returned nil order")
			}

			// requested count matches
			if order.RequestedItemCount != tt.orderRequest.ItemCount {
				t.Errorf("RequestedItemCount = %d, want %d", order.RequestedItemCount, tt.orderRequest.ItemCount)
			}

			// status is set to new
			if order.Status != models.OrderStatusNew {
				t.Errorf("Status = %s, want %s", order.Status, models.OrderStatusNew)
			}

			// shipped item count equal or greater than requested count (overshooting)
			if order.ShippedItemCount < tt.orderRequest.ItemCount {
				t.Errorf("ShippedItemCount = %d, has to be at least %d", order.ShippedItemCount, tt.orderRequest.ItemCount)
			}

			// packs map is not empty
			if order.Packs == nil {
				t.Error("packs map is nil")
			}
			if len(order.Packs) == 0 {
				t.Error("packs map is empty")
			}

			// Verify total packs calculation matches shipped items
			totalItems := 0
			for pack, count := range order.Packs {
				totalItems += int(pack) * count
			}
			if totalItems != order.ShippedItemCount {
				t.Errorf("total items from packs = %d, want %d", totalItems, order.ShippedItemCount)
			}
		})
	}
}

func TestService_CreateOrder_ErrorCases(t *testing.T) {
	tests := []struct {
		name         string
		maxCount     int
		orderRequest models.OrderRequest
		packs        []models.Pack
		expectedErr  error
	}{
		{
			name:         "zero item count",
			maxCount:     1000000000,
			orderRequest: models.OrderRequest{ItemCount: 0},
			packs:        []models.Pack{250, 500, 1000, 2000, 5000},
			expectedErr:  InvalidOrderItemCountError,
		},
		{
			name:         "negative item count",
			maxCount:     1000000000,
			orderRequest: models.OrderRequest{ItemCount: -1},
			packs:        []models.Pack{250, 500, 1000, 2000, 5000},
			expectedErr:  InvalidOrderItemCountError,
		},
		{
			name:         "item count exceeds maximum",
			maxCount:     1000,
			orderRequest: models.OrderRequest{ItemCount: 1001},
			packs:        []models.Pack{250, 500, 1000, 2000, 5000},
			expectedErr:  InvalidOrderItemCountError,
		},
		{
			name:         "no available packs",
			maxCount:     1000000000,
			orderRequest: models.OrderRequest{ItemCount: 1},
			packs:        []models.Pack{},
			expectedErr:  OrderCalculationError,
		},
		{
			name:         "nil available packs",
			maxCount:     1000000000,
			orderRequest: models.OrderRequest{ItemCount: 1},
			packs:        nil,
			expectedErr:  OrderCalculationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.maxCount)
			order, err := service.CreateOrder(tt.orderRequest, tt.packs)

			// has to error
			if err == nil {
				t.Fatalf("CreateOrder() expected error but got order: %+v", order)
			}

			// order is nil
			if order != nil {
				t.Errorf("CreateOrder() expected nil order but got: %+v", order)
			}

			// error is of expected type
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("CreateOrder() error = %v, want error type %v", err, tt.expectedErr)
			}
		})
	}
}
