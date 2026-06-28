package repository

import (
	"errors"

	"github.com/mkamrul9/spotsync-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReservationRepository interface {
	CreateReservationAtomic(reservation *models.Reservation) error
	GetReservationsByUser(userID uint) ([]models.Reservation, error)
	GetReservationByID(id uint) (*models.Reservation, error)
	UpdateReservationStatus(reservation *models.Reservation, status string) error
	GetAllReservations() ([]models.Reservation, error)
	GetActiveReservationCountByZone(zoneID uint) (int64, error)
	// GetActiveCountsPerZone returns a map[zoneID]activeCount in a single query.
	GetActiveCountsPerZone() (map[uint]int64, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db}
}

// CreateReservationAtomic implements Row-Level Locking (FOR UPDATE)
func (r *reservationRepository) CreateReservationAtomic(reservation *models.Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone

		// 1. Lock the row! (SELECT ... FOR UPDATE)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, reservation.ZoneID).Error; err != nil {
			return err
		}

		// 2. Count current 'active' reservations for this zone
		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", reservation.ZoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Check capacity constraints
		if int(activeCount) >= zone.TotalCapacity {
			return errors.New("zone is at maximum capacity")
		}

		// 4. Create the reservation
		if err := tx.Create(reservation).Error; err != nil {
			return err
		}

		// 5. Transaction commits automatically if we return nil
		return nil
	})
}

func (r *reservationRepository) GetReservationsByUser(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	// Using Preload to fetch associated Zone details as required by the assignment
	err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) GetReservationByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.First(&reservation, id).Error
	return &reservation, err
}

func (r *reservationRepository) UpdateReservationStatus(reservation *models.Reservation, status string) error {
	reservation.Status = status
	return r.db.Save(reservation).Error
}

func (r *reservationRepository) GetAllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("User").Preload("Zone").Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) GetActiveReservationCountByZone(zoneID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count).Error
	return count, err
}

// GetActiveCountsPerZone fetches all active reservation counts grouped by zone_id in a
// single SQL query (SELECT zone_id, COUNT(*) … GROUP BY zone_id), eliminating N+1 queries.
func (r *reservationRepository) GetActiveCountsPerZone() (map[uint]int64, error) {
	type result struct {
		ZoneID uint
		Count  int64
	}

	var rows []result
	err := r.db.Model(&models.Reservation{}).
		Select("zone_id, COUNT(*) as count").
		Where("status = ?", "active").
		Group("zone_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[uint]int64, len(rows))
	for _, row := range rows {
		counts[row.ZoneID] = row.Count
	}
	return counts, nil
}
