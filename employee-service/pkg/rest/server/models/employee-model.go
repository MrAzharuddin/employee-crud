package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model	
	Name string `json:"name,omitempty"`

	Position string `json:"position,omitempty"`

	Salary float64 `json:"salary,omitempty"`
}
