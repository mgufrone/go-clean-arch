package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"recipes/domains/shared"
	"recipes/domains/stats"
	"recipes/domains/user"
)

type LikeModel struct {
	gorm.Model
	UserID uint `gorm:"index" json:"user_id"`
	Reference string `gorm:"index" json:"reference"`
	ReferenceID uint64 `gorm:"index" json:"reference_id"`
}

func (v *LikeModel) Transform() *stats.Like {
	return &stats.Like{
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
func ParseLike(v *stats.Like) (res *LikeModel)  {
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
