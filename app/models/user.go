package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"recipes/domains/shared"
	"recipes/domains/user"
)

type UserModel struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `gorm:"unique" json:"email"`
	Password  string `json:"password"`
	Roles     []*RoleModel `gorm:"many2many:user_roles" json:"roles"`
	IsActive  bool `json:"is_active"`
}

func (u *UserModel) Transform() (usr *user.User) {
	if u == nil {
		return nil
	}
	usr = &user.User{
		Model:     shared.Model{
			ID: u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		IsActive:  u.IsActive,
	}
	if len(u.Roles) > 0 {
		for _, rl := range u.Roles {
			usr.Roles = append(usr.Roles, rl.Transform())
		}
	}
	return
}

func ParseUser(usr *user.User) (res *UserModel) {
	if usr == nil {
		return
	}
	b, _ := json.Marshal(usr)
	_ = json.Unmarshal(b, &res)
	return
}
