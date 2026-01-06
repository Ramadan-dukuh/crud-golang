package model

import "time"

type User struct {
	Id         int    `json:"id,omitempty"`
	Nama       string `json:"nama,omitempty"`
	Email      string `json:"email,omitempty"`
	Created_at time.Time `json:"created_at,omitempty"`
	Updated_at time.Time `json:"updated_at,omitempty"`
}