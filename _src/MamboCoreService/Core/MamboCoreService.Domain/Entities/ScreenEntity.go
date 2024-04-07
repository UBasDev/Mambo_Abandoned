package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScreenEntity struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:Id"`
	Name        string    `gorm:"unique;not null;type:string;size:50;column:Name"`
	Value       string    `gorm:"unique;not null;type:string;size:50;column:Value"`
	OrderNumber uint16    `gorm:"unique;not null;type:smallint;column:OrderNumber"`
	//Description string    `gorm:"type:string;size:100;column:Description"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Users     []UserEntity   `gorm:"many2many:UserScreen;foreignKey:Id;joinForeignKey:ScreenId;References:Id;joinReferences:UserId"`
	Roles     []RoleEntity   `gorm:"many2many:ScreenRole;foreignKey:Id;joinForeignKey:ScreenId;References:Id;joinReferences:RoleId"`
}

func (ScreenEntity) TableName() string {

	return "Screens"
}
