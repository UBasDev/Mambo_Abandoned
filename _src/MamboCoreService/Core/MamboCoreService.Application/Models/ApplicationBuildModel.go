package models

type ApplicationBuildModel struct {
	ApiAdress         string
	ReadTimout        uint16
	ReadHeaderTimeout uint16
	WriteTimeout      uint16
	IdleTimeout       uint16
	Environment       string
}
