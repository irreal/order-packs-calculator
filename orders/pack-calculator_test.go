package orders

import (
	"testing"

	"github.com/irreal/order-packs/models"
)

func TestCalculatePack_Functionality(t *testing.T) {
	tests := []struct {
		name              string
		availablePacks    []models.Pack
		requestedCount    int
		expectedPacks     map[models.Pack]int
		expectedItemCount int
		expectedPackCount int
	}{
		{
			name:              "single pack exact match",
			availablePacks:    []models.Pack{5},
			requestedCount:    5,
			expectedPacks:     map[models.Pack]int{5: 1},
			expectedItemCount: 5,
			expectedPackCount: 1,
		},
		{
			name:              "single pack with overshoot",
			availablePacks:    []models.Pack{5},
			requestedCount:    3,
			expectedPacks:     map[models.Pack]int{5: 1},
			expectedItemCount: 5,
			expectedPackCount: 1,
		},
		{
			name:              "multiple same packs",
			availablePacks:    []models.Pack{5},
			requestedCount:    12,
			expectedPacks:     map[models.Pack]int{5: 3},
			expectedItemCount: 15,
			expectedPackCount: 3,
		},
		// examples from requirements
		{
			name:              "order 1 item, choose smallest pack (250)",
			availablePacks:    []models.Pack{250, 500, 1000, 2000, 5000},
			requestedCount:    1,
			expectedPacks:     map[models.Pack]int{250: 1},
			expectedItemCount: 250,
			expectedPackCount: 1,
		},
		{
			name:              "order 250 items, exact match with 1x250",
			availablePacks:    []models.Pack{250, 500, 1000, 2000, 5000},
			requestedCount:    250,
			expectedPacks:     map[models.Pack]int{250: 1},
			expectedItemCount: 250,
			expectedPackCount: 1,
		},
		{
			name:              "order 251 items, prefer 1x500 over 2x250",
			availablePacks:    []models.Pack{250, 500, 1000, 2000, 5000},
			requestedCount:    251,
			expectedPacks:     map[models.Pack]int{500: 1},
			expectedItemCount: 500,
			expectedPackCount: 1,
		},
		{
			name:              "order 501 items, prefer 1x1000 over 2x500 or 3x250",
			availablePacks:    []models.Pack{250, 500, 1000, 2000, 5000},
			requestedCount:    501,
			expectedPacks:     map[models.Pack]int{250: 1, 500: 1},
			expectedItemCount: 750,
			expectedPackCount: 2,
		},
		{
			name:              "order 12001 items, optimal is 2x5000 + 1x2000 + 1x250 over 3x5000",
			availablePacks:    []models.Pack{250, 500, 1000, 2000, 5000},
			requestedCount:    12001,
			expectedPacks:     map[models.Pack]int{5000: 2, 2000: 1, 250: 1},
			expectedItemCount: 12250,
			expectedPackCount: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculatePack(tt.availablePacks, tt.requestedCount)
			if err != nil {
				t.Fatalf("CalculatePack() error = %v", err)
			}

			if result.TotalItems != tt.expectedItemCount {
				t.Errorf("TotalItems = %d, want %d", result.TotalItems, tt.expectedItemCount)
			}

			if result.TotalPacks != tt.expectedPackCount {
				t.Errorf("TotalPacks = %d, want %d", result.TotalPacks, tt.expectedPackCount)
			}

			if len(result.Packs) != len(tt.expectedPacks) {
				t.Errorf("count of pack types used = %d, want %d", len(result.Packs), len(tt.expectedPacks))
			}

			for pack, count := range tt.expectedPacks {
				if result.Packs[pack] != count {
					t.Errorf("pack %d count = %d, want %d", pack, result.Packs[pack], count)
				}
			}
		})
	}
}

func TestCalculatePack_ErrorCases(t *testing.T) {
	tests := []struct {
		name           string
		availablePacks []models.Pack
		requestedCount int
		expectedError  string
	}{
		{
			name:           "zero requested count",
			availablePacks: []models.Pack{5},
			requestedCount: 0,
			expectedError:  "requested count must be greater than 0",
		},
		{
			name:           "negative requested count",
			availablePacks: []models.Pack{5},
			requestedCount: -1,
			expectedError:  "requested count must be greater than 0",
		},
		{
			name:           "empty available packs",
			availablePacks: []models.Pack{},
			requestedCount: 5,
			expectedError:  "no packs available to fulfill the order",
		},
		{
			name:           "nil available packs",
			availablePacks: nil,
			requestedCount: 5,
			expectedError:  "no packs available to fulfill the order",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculatePack(tt.availablePacks, tt.requestedCount)
			if err == nil {
				t.Fatalf("CalculatePack() expected error but got result: %+v", result)
			}

			if err.Error() != tt.expectedError {
				t.Errorf("CalculatePack() error = %v, want %v", err.Error(), tt.expectedError)
			}
		})
	}
}
