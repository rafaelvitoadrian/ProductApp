package model

import "time"

type Product struct {
	Id          int       `json:"id"`
	Id_store    int       `json:"id_store"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}
