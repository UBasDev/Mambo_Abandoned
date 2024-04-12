package models

import enums "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Enums"

type ApplicationContainerModel struct {
	ApiAdress         string
	ReadTimout        uint16
	ReadHeaderTimeout uint16
	WriteTimeout      uint16
	IdleTimeout       uint16
	Environment       enums.Environment
}
