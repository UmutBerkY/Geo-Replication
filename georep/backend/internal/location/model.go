package location

import "time"

// Location represents a stored geographical point.
type Location struct {
	ID        int64     `json:"id"`
	City      string    `json:"city"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateLocationInput is the payload for creating a location.
type CreateLocationInput struct {
	City string  `json:"city" binding:"required"`
	Lat  float64 `json:"lat" binding:"required"`
	Lon  float64 `json:"lon" binding:"required"`
}



