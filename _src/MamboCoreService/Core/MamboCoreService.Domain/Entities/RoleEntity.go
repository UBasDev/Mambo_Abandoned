package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleEntity struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:Id"`
	Name      string    `gorm:"unique;not null;type:string;size:50;column:Name"`
	Value     string    `gorm:"unique;not null;type:string;size:50;column:Value"`
	ShortCode string    `gorm:"unique;not null;type:string;size:10;column:ShortCode"`
	Level     uint16    `gorm:"unique;not null;type:smallint;column:Level"`
	//Description string    `gorm:"type:string;size:100;column:Description"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Users     []UserEntity   `gorm:"foreignKey:RoleId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Screens   []ScreenEntity `gorm:"many2many:ScreenRole;foreignKey:Id;joinForeignKey:RoleId;References:Id;joinReferences:ScreenId"`
}

func (RoleEntity) TableName() string {
	return "Roles"
}
