package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileEntity struct {
	Id        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:Id"`
	Firstname string    `gorm:"not null;type:string;size:50;column:Firstname"`
	Lastname  string    `gorm:"not null;type:string;size:50;column:Lastname"`
	Age       uint8     `gorm:"type:smallint"`
	Gender    string    `gorm:"type:varchar(10)"`
	BirthDate time.Time `gorm:"column:BirthDate"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserId    uuid.UUID      `gorm:"type:uuid;column:UserId"`
	CompanyId uuid.UUID      `gorm:"type:uuid;column:CompanyId"`
}

func (ProfileEntity) TableName() string {
	return "Profiles"
}
