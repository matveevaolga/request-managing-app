package domain

import (
	"time"
)

type ProjectType struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *ProjectType) Validate() error {
	if len(p.Name) < 2 || len(p.Name) > 100 {
		return ErrInvalidProjectTypeName
	}
	return nil
}
