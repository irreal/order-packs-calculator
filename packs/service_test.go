package packs

import (
	"errors"
	"reflect"
	"testing"

	"github.com/irreal/order-packs/models"
)

type MockPackRepository struct {
	packs          models.Packs
	getPacksError  error
	savePacksError error
}

func NewMockPackRepository() *MockPackRepository {
	return &MockPackRepository{
		packs: make(models.Packs, 0),
	}
}

func (m *MockPackRepository) GetPacks() (models.Packs, error) {
	if m.getPacksError != nil {
		return nil, m.getPacksError
	}
	return m.packs, nil
}

func (m *MockPackRepository) SavePacks(packs models.Packs) error {
	if m.savePacksError != nil {
		return m.savePacksError
	}
	m.packs = packs
	return nil
}

func (m *MockPackRepository) SetPacks(packs models.Packs) {
	m.packs = packs
}

func (m *MockPackRepository) SetGetPacksError(err error) {
	m.getPacksError = err
}

func (m *MockPackRepository) SetSavePacksError(err error) {
	m.savePacksError = err
}

func (m *MockPackRepository) Reset() {
	m.packs = make(models.Packs, 0)
	m.getPacksError = nil
	m.savePacksError = nil
}

func TestService_GetPacks(t *testing.T) {
	tests := []struct {
		name          string
		existingPacks models.Packs
		expectedPacks models.Packs
	}{
		{
			name:          "empty packs",
			existingPacks: models.Packs{},
			expectedPacks: models.Packs{},
		},
		{
			name:          "single pack",
			existingPacks: models.Packs{250},
			expectedPacks: models.Packs{250},
		},
		{
			name:          "multiple packs",
			existingPacks: models.Packs{250, 500, 1000, 2000, 5000},
			expectedPacks: models.Packs{250, 500, 1000, 2000, 5000},
		},
		{
			name:          "unsorted packs",
			existingPacks: models.Packs{1000, 250, 5000, 500, 2000},
			expectedPacks: models.Packs{1000, 250, 5000, 500, 2000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockPackRepository()
			mockRepo.SetPacks(tt.existingPacks)
			service := NewService(mockRepo)

			packs, err := service.GetPacks()

			if err != nil {
				t.Fatalf("GetPacks() unexpected error = %v", err)
			}

			if !reflect.DeepEqual(packs, tt.expectedPacks) {
				t.Errorf("GetPacks() = %v, want %v", packs, tt.expectedPacks)
			}
		})
	}
}

func TestService_GetPacks_RepositoryError(t *testing.T) {
	mockRepo := NewMockPackRepository()
	mockRepo.SetGetPacksError(errors.New("database connection failed"))
	service := NewService(mockRepo)

	packs, err := service.GetPacks()

	// return error when repository fails
	if err == nil {
		t.Fatal("GetPacks() expected error when repository fails")
	}

	// return nil packs when repository fails
	if packs != nil {
		t.Errorf("GetPacks() expected nil packs when repository fails, got: %+v", packs)
	}

	// error should be the repository error
	expectedError := "database connection failed"
	if err.Error() != expectedError {
		t.Errorf("GetPacks() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestService_SavePacks(t *testing.T) {
	tests := []struct {
		name        string
		packsToSave models.Packs
	}{
		{
			name:        "empty packs",
			packsToSave: models.Packs{},
		},
		{
			name:        "single pack",
			packsToSave: models.Packs{250},
		},
		{
			name:        "multiple packs",
			packsToSave: models.Packs{250, 500, 1000, 2000, 5000},
		},
		{
			name:        "large pack sizes",
			packsToSave: models.Packs{10000, 50000, 100000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockPackRepository()
			service := NewService(mockRepo)

			err := service.SavePacks(tt.packsToSave)

			if err != nil {
				t.Fatalf("SavePacks() unexpected error = %v", err)
			}

			// packs were saved to the repository
			savedPacks, _ := mockRepo.GetPacks()
			if !reflect.DeepEqual(savedPacks, tt.packsToSave) {
				t.Errorf("SavePacks() saved packs = %v, want %v", savedPacks, tt.packsToSave)
			}
		})
	}
}

func TestService_SavePacks_RepositoryError(t *testing.T) {
	mockRepo := NewMockPackRepository()
	mockRepo.SetSavePacksError(errors.New("disk full"))
	service := NewService(mockRepo)

	packsToSave := models.Packs{250, 500, 1000}
	err := service.SavePacks(packsToSave)

	// return error when repository fails
	if err == nil {
		t.Fatal("SavePacks() expected error when repository fails")
	}

	// error should be the repository error
	expectedError := "disk full"
	if err.Error() != expectedError {
		t.Errorf("SavePacks() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestService_SaveAndRetrievePacks(t *testing.T) {
	mockRepo := NewMockPackRepository()
	service := NewService(mockRepo)

	originalPacks := models.Packs{100, 250, 500, 1000}

	// save packs
	err := service.SavePacks(originalPacks)
	if err != nil {
		t.Fatalf("SavePacks() unexpected error = %v", err)
	}

	// retrieve packs
	retrievedPacks, err := service.GetPacks()
	if err != nil {
		t.Fatalf("GetPacks() unexpected error = %v", err)
	}

	// match
	if !reflect.DeepEqual(retrievedPacks, originalPacks) {
		t.Errorf("Retrieved packs = %v, want %v", retrievedPacks, originalPacks)
	}
}
