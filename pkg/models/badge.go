package models

type Badge struct {
	// A badge is a small icon that appears next to a user's name.
	Model
	Id              string `json:"id" gorm:"primaryKey"`
	BackgroundColor string `json:"backgroundColor"`
	TextColor       string `json:"textColor"`
}
