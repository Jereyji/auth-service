package models

import "time"

type RegisterEvent struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Link      string    `json:"link"`
}

type LoginEvent struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
}
