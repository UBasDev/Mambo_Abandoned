package entities

import (
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserEntity interface {
	GetId() uuid.UUID
	GetUsername() string
	SetUsername(username string)
	GetEmail() string
	SetEmail(email string)
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDeletedAt() gorm.DeletedAt
	GetRoleId() *uuid.UUID
	SetRoleId(*uuid.UUID)
}
type UserEntity struct {
	Id            uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:Id"`
	Username      string    `gorm:"unique;not null;type:string;size:50;column:Username"`
	Email         string    `gorm:"unique;not null;type:string;size:50;column:Email"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt  `gorm:"index"`
	RoleId        *uuid.UUID      `gorm:"type:uuid;column:RoleId"`
	Screens       []*ScreenEntity `gorm:"many2many:UserScreen;foreignKey:Id;joinForeignKey:UserId;References:Id;joinReferences:ScreenId"`
	ProfileEntity *ProfileEntity  `gorm:"foreignKey:UserId;references:Id"`
}

func BuildNewUserEntity(username string, email string) IUserEntity {
	return &UserEntity{
		Username: username,
		Email:    email,
	}
}
func (u *UserEntity) GetId() uuid.UUID {
	return u.Id
}
func (u *UserEntity) GetUsername() string {
	return u.Username
}
func (u *UserEntity) SetUsername(username string) {
	u.Username = username
}
func (u *UserEntity) GetEmail() string {
	return u.Email
}
func (u *UserEntity) SetEmail(email string) {
	u.Email = email
}
func (u *UserEntity) GetCreatedAt() time.Time {
	return u.CreatedAt
}
func (u *UserEntity) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}
func (u *UserEntity) GetDeletedAt() gorm.DeletedAt {
	return u.DeletedAt
}
func (u *UserEntity) GetRoleId() *uuid.UUID {
	return u.RoleId
}
func (u *UserEntity) SetRoleId(roleId *uuid.UUID) {
	u.RoleId = roleId
}
func (UserEntity) TableName() string {
	return "Users"
}
