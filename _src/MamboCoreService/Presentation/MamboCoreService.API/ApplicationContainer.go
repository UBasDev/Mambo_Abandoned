package main

import (
	"net/http"

	enums "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Enums"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type applicationContainerModel struct {
	ApiAdress           string
	ReadTimout          uint16
	ReadHeaderTimeout   uint16
	WriteTimeout        uint16
	IdleTimeout         uint16
	Environment         enums.Environment
	ZapLogger           *zap.Logger
	PostgreSqlDbContext *gorm.DB
	Router              *http.ServeMux
	Server              http.Server
}

func createNewApplicationContainerModel(apiAdress string, readTimout uint16, readHeaderTimeout uint16, writeTimeout uint16, idleTimeout uint16, environment enums.Environment) applicationContainerModel {
	return applicationContainerModel{
		ApiAdress:         apiAdress,
		ReadTimout:        readTimout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		Environment:       environment,
	}
}
