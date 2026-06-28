package repository

import (
	"github.com/mkamrul9/spotsync-api/models"
	"gorm.io/gorm"
)

type ZoneRepository interface {
	CreateZone(zone *models.ParkingZone) error
	GetAllZones() ([]models.ParkingZone, error)
	GetZoneByID(id uint) (*models.ParkingZone, error)
	UpdateZone(zone *models.ParkingZone) error
	DeleteZone(id uint) error
}

type zoneRepository struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepository{db}
}

func (r *zoneRepository) CreateZone(zone *models.ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *zoneRepository) GetAllZones() ([]models.ParkingZone, error) {
	var zones []models.ParkingZone
	err := r.db.Find(&zones).Error
	return zones, err
}

func (r *zoneRepository) GetZoneByID(id uint) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	err := r.db.First(&zone, id).Error
	return &zone, err
}

// UpdateZone persists field-level changes to an existing zone using GORM's Save.
func (r *zoneRepository) UpdateZone(zone *models.ParkingZone) error {
	return r.db.Save(zone).Error
}

// DeleteZone performs a hard delete of the zone by its primary key.
func (r *zoneRepository) DeleteZone(id uint) error {
	return r.db.Delete(&models.ParkingZone{}, id).Error
}
