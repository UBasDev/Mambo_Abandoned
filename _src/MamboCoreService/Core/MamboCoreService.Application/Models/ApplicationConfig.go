package models

type ApplicationConfig struct {
	AppSettings appSettings `json:"AppSettings"`
}
type appSettings struct {
	Server                    server                    `json:"Server"`
	DatabaseConnectionStrings databaseConnectionStrings `json:"DatabaseConnectionStrings"`
	LogFilePath               string                    `json:"LogFilePath"`
}
type server struct {
	Host              string `json:"Host"`
	Port              string `json:"Port"`
	ReadTimeout       uint16 `json:"ReadTimeout"`
	ReadHeaderTimeout uint16 `json:"ReadHeaderTimeout"`
	WriteTimeout      uint16 `json:"WriteTimeout"`
	IdleTimeout       uint16 `json:"IdleTimeout"`
}
type databaseConnectionStrings struct {
	PostgreSqlDbUrl string `json:"PostgreSqlDbUrl"`
}
