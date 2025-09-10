package orders

import (
	"errors"
	"strings"
	"testing"

	"github.com/irreal/order-packs/models"
)

type MockOrderRepository struct {
	savedOrders    []*models.Order
	saveOrderError error
	getLast10Error error
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		savedOrders: make([]*models.Order, 0),
	}
}

func (m *MockOrderRepository) SaveOrder(order *models.Order) error {
	if m.saveOrderError != nil {
		return m.saveOrderError
	}
	m.savedOrders = append(m.savedOrders, order)
	return nil
}

func (m *MockOrderRepository) GetLast10Orders() ([]*models.Order, error) {
	if m.getLast10Error != nil {
		return nil, m.getLast10Error
	}

	start := 0
	if len(m.savedOrders) > 10 {
		start = len(m.savedOrders) - 10
	}

	result := make([]*models.Order, len(m.savedOrders)-start)
	copy(result, m.savedOrders[start:])
	return result, nil
}

func (m *MockOrderRepository) GetSavedOrders() []*models.Order {
	return m.savedOrders
}

func (m *MockOrderRepository) SetSaveOrderError(err error) {
	m.saveOrderError = err
}

func (m *MockOrderRepository) SetGetLast10Error(err error) {
	m.getLast10Error = err
}

func (m *MockOrderRepository) Reset() {
	m.savedOrders = make([]*models.Order, 0)
	m.saveOrderError = nil
	m.getLast10Error = nil
}

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
			mockRepo := NewMockOrderRepository()
			service := NewService(tt.maxCount, mockRepo)
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

			// total packs calculation matches shipped items
			totalItems := 0
			for pack, count := range order.Packs {
				totalItems += int(pack) * count
			}
			if totalItems != order.ShippedItemCount {
				t.Errorf("total items from packs = %d, want %d", totalItems, order.ShippedItemCount)
			}

			// repository called
			savedOrders := mockRepo.GetSavedOrders()
			if len(savedOrders) != 1 {
				t.Errorf("expected 1 order to be saved, got %d", len(savedOrders))
			} else {
				savedOrder := savedOrders[0]
				if savedOrder.RequestedItemCount != order.RequestedItemCount {
					t.Errorf("saved order RequestedItemCount = %d, want %d", savedOrder.RequestedItemCount, order.RequestedItemCount)
				}
				if savedOrder.ShippedItemCount != order.ShippedItemCount {
					t.Errorf("saved order ShippedItemCount = %d, want %d", savedOrder.ShippedItemCount, order.ShippedItemCount)
				}
				if savedOrder.Status != order.Status {
					t.Errorf("saved order Status = %s, want %s", savedOrder.Status, order.Status)
				}
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
			mockRepo := NewMockOrderRepository()
			service := NewService(tt.maxCount, mockRepo)
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

			// no orders were saved when errors occur
			savedOrders := mockRepo.GetSavedOrders()
			if len(savedOrders) != 0 {
				t.Errorf("expected no orders to be saved on error, but got %d saved orders", len(savedOrders))
			}
		})
	}
}

func TestService_CreateOrder_RepositoryError(t *testing.T) {
	mockRepo := NewMockOrderRepository()
	mockRepo.SetSaveOrderError(errors.New("database connection failed"))
	service := NewService(1000000000, mockRepo)

	orderRequest := models.OrderRequest{ItemCount: 1}
	packs := []models.Pack{250, 500, 1000}

	order, err := service.CreateOrder(orderRequest, packs)

	// return error when repository fails
	if err == nil {
		t.Fatal("CreateOrder() expected error when repository fails")
	}

	// return nil order when repository fails
	if order != nil {
		t.Errorf("CreateOrder() expected nil order when repository fails, got: %+v", order)
	}

	// should mention saving failure
	expectedMsg := "failed to save order"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("CreateOrder() error should contain '%s', got: %v", expectedMsg, err)
	}
}

func TestService_GetLast10Orders(t *testing.T) {
	mockRepo := NewMockOrderRepository()
	service := NewService(1000000000, mockRepo)

	// empty repo
	orders, err := service.GetLast10Orders()
	if err != nil {
		t.Fatalf("GetLast10Orders() unexpected error: %v", err)
	}
	if len(orders) != 0 {
		t.Errorf("GetLast10Orders() expected 0 orders, got %d", len(orders))
	}

	// add mock orders
	testOrders := []*models.Order{
		{RequestedItemCount: 1, ShippedItemCount: 250, Status: models.OrderStatusNew},
		{RequestedItemCount: 2, ShippedItemCount: 250, Status: models.OrderStatusNew},
		{RequestedItemCount: 3, ShippedItemCount: 250, Status: models.OrderStatusNew},
	}

	for _, order := range testOrders {
		mockRepo.SaveOrder(order)
	}

	// retrieval
	orders, err = service.GetLast10Orders()
	if err != nil {
		t.Fatalf("GetLast10Orders() unexpected error: %v", err)
	}
	if len(orders) != 3 {
		t.Errorf("GetLast10Orders() expected 3 orders, got %d", len(orders))
	}

	// orders match
	for i, order := range orders {
		if order.RequestedItemCount != testOrders[i].RequestedItemCount {
			t.Errorf("GetLast10Orders() order %d RequestedItemCount = %d, want %d",
				i, order.RequestedItemCount, testOrders[i].RequestedItemCount)
		}
	}
}

func TestService_GetLast10Orders_RepositoryError(t *testing.T) {
	mockRepo := NewMockOrderRepository()
	mockRepo.SetGetLast10Error(errors.New("database query failed"))
	service := NewService(1000000000, mockRepo)

	orders, err := service.GetLast10Orders()

	// return error when repository fails
	if err == nil {
		t.Fatal("GetLast10Orders() expected error when repository fails")
	}

	// return nil orders when repository fails
	if orders != nil {
		t.Errorf("GetLast10Orders() expected nil orders when repository fails, got: %+v", orders)
	}

	// error should be the repository error
	if !errors.Is(err, mockRepo.getLast10Error) {
		t.Errorf("GetLast10Orders() error = %v, want %v", err, mockRepo.getLast10Error)
	}
}
