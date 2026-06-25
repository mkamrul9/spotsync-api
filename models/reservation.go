package models

import (
	"time"
)

type Reservation struct {
	ID           uint        `gorm:"primaryKey"`
	UserID       uint        `gorm:"not null"`
	User         User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ZoneID       uint        `gorm:"not null"`
	Zone         ParkingZone `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LicensePlate string      `gorm:"type:varchar(15);not null"`
	Status       string      `gorm:"type:varchar(20);default:'active'"` // 'active', 'completed', or 'cancelled'
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
