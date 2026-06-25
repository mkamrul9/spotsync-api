package service

import (
	"errors"

	"github.com/mkamrul9/spotsync-api/dto"
	"github.com/mkamrul9/spotsync-api/models"
	"github.com/mkamrul9/spotsync-api/repository"
)

type ReservationService interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.CreateReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.ReservationResponse, error)
	CancelReservation(userID uint, reservationID uint, role string) error
	GetAllSystemReservations() ([]dto.AdminReservationResponse, error)
}

type reservationService struct {
	resRepo repository.ReservationRepository
}

func NewReservationService(resRepo repository.ReservationRepository) ReservationService {
	return &reservationService{resRepo}
}

func (s *reservationService) CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.CreateReservationResponse, error) {
	res := models.Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       "active",
	}

	// Calls the strict atomic function with Row-Level Locking
	if err := s.resRepo.CreateReservationAtomic(&res); err != nil {
		return nil, err
	}

	return &dto.CreateReservationResponse{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		CreatedAt:    res.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    res.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.ReservationResponse, error) {
	reservations, err := s.resRepo.GetReservationsByUser(userID)
	if err != nil {
		return nil, err
	}

	var response []dto.ReservationResponse
	for _, r := range reservations {
		response = append(response, dto.ReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneSummary{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return response, nil
}

func (s *reservationService) CancelReservation(userID uint, reservationID uint, role string) error {
	res, err := s.resRepo.GetReservationByID(reservationID)
	if err != nil {
		return errors.New("reservation not found")
	}

	// 403 Forbidden check: Drivers can only cancel their own
	if role != "admin" && res.UserID != userID {
		return errors.New("forbidden: you can only cancel your own reservations")
	}

	return s.resRepo.UpdateReservationStatus(res, "cancelled")
}

func (s *reservationService) GetAllSystemReservations() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.resRepo.GetAllReservations()
	if err != nil {
		return nil, err
	}

	var response []dto.AdminReservationResponse
	for _, r := range reservations {
		response = append(response, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			User: dto.UserSummary{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
			},
			Zone: dto.ZoneSummary{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return response, nil
}
