package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/amend-parking-backend/internal/config"
	"github.com/amend-parking-backend/internal/models"
	"github.com/amend-parking-backend/internal/repository"
	"github.com/google/uuid"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCountOfFreeSpaces(ctx context.Context) (int, error) {
	occupiedCount, err := s.repo.GetCountOfOccupiedSpaces(ctx)
	if err != nil {
		return 0, err
	}
	return config.Settings.ParkingSlotsCount - int(occupiedCount), nil
}

func (s *Service) GetOccupiedSpaces(ctx context.Context) ([]models.ParkingSpaceLog, error) {
	return s.repo.GetOccupiedSpaces(ctx)
}

func (s *Service) AddParkingSpaceLog(ctx context.Context, firstName, lastName, carMake, licensePlate string) (*models.ParkingSpaceLog, error) {
	occupiedSpaces, err := s.repo.GetOccupiedSpaces(ctx)
	if err != nil {
		return nil, err
	}

	if len(occupiedSpaces) >= config.Settings.ParkingSlotsCount {
		return nil, fmt.Errorf("no free parking spaces available")
	}

	occupiedPlaceNumbers := make(map[int]bool)
	for _, space := range occupiedSpaces {
		occupiedPlaceNumbers[space.PlaceNumber] = true
	}

	var availablePlaces []int
	for i := 1; i <= config.Settings.ParkingSlotsCount; i++ {
		if !occupiedPlaceNumbers[i] {
			availablePlaces = append(availablePlaces, i)
		}
	}

	if len(availablePlaces) == 0 {
		return nil, fmt.Errorf("no free parking spaces available")
	}

	rand.Seed(time.Now().UnixNano())
	selectedPlace := availablePlaces[rand.Intn(len(availablePlaces))]

	parkingSpaceLog := &models.ParkingSpaceLog{
		LogID:        uuid.New().String(),
		PlaceNumber:  selectedPlace,
		FirstName:    firstName,
		LastName:     lastName,
		CarMake:      carMake,
		LicensePlate: licensePlate,
		CreatedAt:    time.Now().UTC(),
		IsActive:     true,
	}

	err = s.repo.AddParkingSpaceLog(ctx, parkingSpaceLog)
	if err != nil {
		return nil, err
	}

	return parkingSpaceLog, nil
}

func (s *Service) FreeUpParkingSpace(ctx context.Context, placeNumber int) (*models.ParkingSpaceLog, error) {
	parkingSpaceLog, err := s.repo.GetParkingSpaceLogByPlaceNumber(ctx, placeNumber)
	if err != nil {
		return nil, err
	}

	if parkingSpaceLog == nil {
		return nil, fmt.Errorf("parking space is already free")
	}

	now := time.Now().UTC()
	parkingSpaceLog.IsActive = false
	parkingSpaceLog.FreeUpTime = &now

	err = s.repo.UpdateParkingSpaceLog(ctx, parkingSpaceLog)
	if err != nil {
		return nil, err
	}

	return parkingSpaceLog, nil
}

func (s *Service) GetParkingSpaceLogsByFirstNameAndLastName(ctx context.Context, firstName, lastName string) ([]models.ParkingSpaceLog, error) {
	return s.repo.GetParkingSpaceLogsByFirstNameAndLastName(ctx, firstName, lastName)
}
