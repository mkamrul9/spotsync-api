package dto

// CreateReservationRequest defines the payload for booking a spot
type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required,gt=0"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

// ReservationResponse handles the output structure
type ReservationResponse struct {
	ID           uint        `json:"id"`
	LicensePlate string      `json:"license_plate"`
	Status       string      `json:"status"`
	Zone         ZoneSummary `json:"zone"`
	CreatedAt    string      `json:"created_at"`
}

type ZoneSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateReservationResponse struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	ZoneID       uint   `json:"zone_id"`
	LicensePlate string `json:"license_plate"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type AdminReservationResponse struct {
	ID           uint        `json:"id"`
	LicensePlate string      `json:"license_plate"`
	Status       string      `json:"status"`
	User         UserSummary `json:"user"`
	Zone         ZoneSummary `json:"zone"`
	CreatedAt    string      `json:"created_at"`
}
