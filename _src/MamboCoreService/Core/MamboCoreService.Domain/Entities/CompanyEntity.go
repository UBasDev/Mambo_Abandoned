package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyEntity struct {
	Id   uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:Id"`
	Name string    `gorm:"unique;not null;type:string;size:50;column:Name"`
	//Description string    `gorm:"type:string;size:100;column:Description"`
	//Adress    string `gorm:"type:string;size:100;column:Adress"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt  `gorm:"index"`
	Profiles  []ProfileEntity `gorm:"foreignKey:CompanyId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (CompanyEntity) TableName() string {

	return "Companies"
}
