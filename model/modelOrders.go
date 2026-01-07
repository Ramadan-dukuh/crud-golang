package model

import "time"

type Orders struct {
	Id          int       `json:"id,omitempty"`
	User_id     int       `json:"user_id,omitempty"`
	Total_harga int       `json:"total_harga,omitempty"`
	Status      string    `json:"status,omitempty"`
	Created_at  time.Time `json:"created_at,omitempty"`
}
