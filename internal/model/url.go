package model

import "time"

type Url struct {
	Short     string    `json:"short_code"`
	Original  string    `json:"original_url"`
	CreatedAt time.Time `json:"created_at"`
	Clicks    int       `json:"clicks"`
}
