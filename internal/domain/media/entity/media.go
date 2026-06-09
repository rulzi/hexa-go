package entity

import "time"

// Media represents the media entity in the domain
type Media struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate validates the media entity
func (m *Media) Validate() error {
	if m.Name == "" {
		return NewNameRequired()
	}
	if m.Path == "" {
		return NewPathRequired()
	}
	return nil
}
