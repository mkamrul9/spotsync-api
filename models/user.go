package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`                          // Stores the bcrypt hashed password
	Role      string `gorm:"type:varchar(20);default:'driver'"` // 'driver' or 'admin'
	CreatedAt time.Time
	UpdatedAt time.Time
}
