package models

import "gorm.io/gorm"

type Place struct {
	gorm.Model
	PlaceID              string  `json:"place_id"`
	Name                 string  `json:"name"`
	FormattedPhoneNumber string  `json:"formatted_phone_number"`
	PriceLevel           int     `json:"price_level"`
	Rating               float64 `json:"rating"`
	Types                string  `json:"types"`
	Vicinity             string  `json:"vicinity"`
	Contacted            bool
	ContactName          string
	Notes                string
}
