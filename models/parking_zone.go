package models

import (
	"time"
)

type ParkingZone struct {
	ID            uint    `gorm:"primaryKey"`
	Name          string  `gorm:"not null"`
	Type          string  `gorm:"type:varchar(50);not null"` // 'general', 'ev_charging', or 'covered'
	TotalCapacity int     `gorm:"not null;check:total_capacity > 0"`
	PricePerHour  float64 `gorm:"type:decimal(10,2);not null;check:price_per_hour > 0"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
