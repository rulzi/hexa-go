package entity

import "time"

// Media represents the media entity in the domain
type Media struct {
	ID        int64
	Name      string
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
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
