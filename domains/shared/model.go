package shared

import "time"

type Model struct {
	ID interface{} `json:",omitempty"`
	CreatedAt time.Time `json:",omitempty"`
	UpdatedAt time.Time `json:",omitempty"`
}
