package service

import (
	"github.com/mkamrul9/spotsync-api/dto"
	"github.com/mkamrul9/spotsync-api/models"
	"github.com/mkamrul9/spotsync-api/repository"
)

type ZoneService interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.CreateZoneResponse, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
	UpdateZone(id uint, req dto.UpdateZoneRequest) (*dto.CreateZoneResponse, error)
	DeleteZone(id uint) error
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
	resRepo  repository.ReservationRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository, resRepo repository.ReservationRepository) ZoneService {
	return &zoneService{zoneRepo, resRepo}
}

func (s *zoneService) CreateZone(req dto.CreateZoneRequest) (*dto.CreateZoneResponse, error) {
	zone := models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zoneRepo.CreateZone(&zone); err != nil {
		return nil, err
	}

	return &dto.CreateZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     zone.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, err := s.zoneRepo.GetAllZones()
	if err != nil {
		return nil, err
	}

	// Fetch all active reservation counts in a single GROUP BY query (avoids N+1)
	activeCounts, err := s.resRepo.GetActiveCountsPerZone()
	if err != nil {
		return nil, err
	}

	var response []dto.ZoneResponse
	for _, z := range zones {
		available := z.TotalCapacity - int(activeCounts[z.ID]) // 0 if zone has no active reservations

		response = append(response, dto.ZoneResponse{
			ID:             z.ID,
			Name:           z.Name,
			Type:           z.Type,
			TotalCapacity:  z.TotalCapacity,
			AvailableSpots: available,
			PricePerHour:   z.PricePerHour,
			CreatedAt:      z.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return response, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	z, err := s.zoneRepo.GetZoneByID(id)
	if err != nil {
		return nil, err
	}

	activeCount, _ := s.resRepo.GetActiveReservationCountByZone(z.ID)
	available := z.TotalCapacity - int(activeCount)

	return &dto.ZoneResponse{
		ID:             z.ID,
		Name:           z.Name,
		Type:           z.Type,
		TotalCapacity:  z.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   z.PricePerHour,
		CreatedAt:      z.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// UpdateZone applies partial updates to an existing zone. Only non-nil fields are changed.
func (s *zoneService) UpdateZone(id uint, req dto.UpdateZoneRequest) (*dto.CreateZoneResponse, error) {
	// 1. Fetch the existing zone
	zone, err := s.zoneRepo.GetZoneByID(id)
	if err != nil {
		return nil, err
	}

	// 2. Apply only the fields that were provided (pointer check)
	if req.Name != nil {
		zone.Name = *req.Name
	}
	if req.Type != nil {
		zone.Type = *req.Type
	}
	if req.TotalCapacity != nil {
		zone.TotalCapacity = *req.TotalCapacity
	}
	if req.PricePerHour != nil {
		zone.PricePerHour = *req.PricePerHour
	}

	// 3. Persist
	if err := s.zoneRepo.UpdateZone(zone); err != nil {
		return nil, err
	}

	return &dto.CreateZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     zone.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// DeleteZone removes the zone by ID.
func (s *zoneService) DeleteZone(id uint) error {
	// Verify zone exists first so we can return a proper 404
	if _, err := s.zoneRepo.GetZoneByID(id); err != nil {
		return err
	}
	return s.zoneRepo.DeleteZone(id)
}
