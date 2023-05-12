package models

import "gorm.io/gorm"

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Zone struct {
	gorm.Model
	ID      uint     `json:"id" gorm:"primary_key"`
	Vertex1 Position `json:"vertex1" gorm:"embedded; embeddedPrefix:vertex1"`
	Vertex2 Position `json:"vertex2" gorm:"embedded; embeddedPrefix:vertex2"`
	Vertex3 Position `json:"vertex3" gorm:"embedded; embeddedPrefix:vertex3"`
	Vertex4 Position `json:"vertex4" gorm:"embedded; embeddedPrefix:vertex4"`
}

type Device struct {
	gorm.Model
	ID        uint     `json:"id" gorm:"primary_key"`
	Position  Position `json:"position" gorm:"embedded"`
	TimeStamp string   `json:"timestamp"`
}
