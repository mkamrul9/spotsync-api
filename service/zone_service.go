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

	var response []dto.ZoneResponse
	for _, z := range zones {
		// Calculate available spots dynamically
		activeCount, _ := s.resRepo.GetActiveReservationCountByZone(z.ID)
		available := z.TotalCapacity - int(activeCount)

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
