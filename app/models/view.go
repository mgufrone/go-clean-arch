package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"recipes/domains/shared"
	"recipes/domains/stats"
	"recipes/domains/user"
)

type ViewModel struct {
	gorm.Model
	UserID uint `gorm:"index" json:"user_id"`
	Reference string `gorm:"index" json:"reference"`
	ReferenceID uint64 `gorm:"index" json:"reference_id"`
}
type ViewModelWithStats struct {
	*ViewModel
	ViewCount int64 `json:"view_count"`
}

func (v *ViewModel) Transform() *stats.View {
	return &stats.View{
		Model:       shared.Model{
			ID: v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		},
		User:        &user.User{
			Model: shared.Model{
				ID: v.UserID,
			},
		},
		Reference:   v.Reference,
		ReferenceID: v.ReferenceID,
	}
}
func ParseView(v *stats.View) (res *ViewModel)  {
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &res); err != nil {
		return
	}
	if v.User != nil && v.User.ID != nil {
		res.UserID = v.User.ID.(uint)
	}
	return
}
