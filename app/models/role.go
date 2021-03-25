package models

import (
	"gorm.io/gorm"
	"recipes/domains/shared"
	"recipes/domains/user"
)

type RoleModel struct {
	gorm.Model
	Name string `gorm:"uniqueIndex,class:FULLTEXT;length:100" json:"name"`
	IsActive *bool `gorm:"index,default:true" json:"is_active"`
	Users []*UserModel `gorm:"many2many:user_roles" json:"users"`
}

func (m *RoleModel) Transform() (role *user.Role) {
	if m == nil {
		return nil
	}
	role = &user.Role{
		Model:    shared.Model{
			ID: m.ID,
		},
		Name:     m.Name,
		IsActive: *m.IsActive,
	}
	if len(m.Users) > 0 {
		for _, m1 := range m.Users {
			role.Users = append(role.Users, m1.Transform())
		}
	}
	return
}

func ParseRole(r *user.Role) (role *RoleModel) {
	if r == nil {
		return
	}
	var id uint
	if r.ID != nil {
		id = r.ID.(uint)
	}
	role = &RoleModel{
		Model:    gorm.Model{
			ID: id,
		},
		Name:     r.Name,
		IsActive: &r.IsActive,
	}
	if len(r.Users) > 0 {
		for _, usr := range r.Users {
			role.Users = append(role.Users, ParseUser(usr))
		}
	}
	return
}


