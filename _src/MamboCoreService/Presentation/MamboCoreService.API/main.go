package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	contexts "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Contexts"
	enums "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Enums"
	models "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Models"
	routes "github.com/UBasDev/Mambo/_src/MamboCoreService/Presentation/MamboCoreService.API/Routes"
	helpers "github.com/UBasDev/Mambo/_src/_helpers"
)

func main() {
	environment := os.Getenv("GOLANG_ENVIRONMENT")
	if environment == "" {
		environment = string(enums.Development)
	}
	applicationConfig := helpers.ConfigFileRead[models.ApplicationConfig](environment)
	applicationBuilder := models.ApplicationBuildModel{
		ApiAdress:         fmt.Sprintf("%s:%s", applicationConfig.AppSettings.Server.Host, applicationConfig.AppSettings.Server.Port),
		ReadTimout:        applicationConfig.AppSettings.Server.ReadTimeout,
		ReadHeaderTimeout: applicationConfig.AppSettings.Server.ReadHeaderTimeout,
		WriteTimeout:      applicationConfig.AppSettings.Server.WriteTimeout,
		IdleTimeout:       applicationConfig.AppSettings.Server.IdleTimeout,
		Environment:       environment,
	}
	postgreSqlDbContext := contexts.InitPostgreSqlDatabaseContext(applicationConfig.AppSettings.DatabaseConnectionStrings.PostgreSqlDbUrl, environment)
	router := routes.CreateRoutes(postgreSqlDbContext)
	server := createServer(router, applicationBuilder.ApiAdress, applicationBuilder.ReadTimout, applicationBuilder.ReadHeaderTimeout, applicationBuilder.WriteTimeout, applicationBuilder.IdleTimeout)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
func createServer(router http.Handler, apiAdress string, readTimeout uint16, readHeaderTimeout uint16, writeTimeout uint16, idleTimeout uint16) http.Server {
	return http.Server{
		Handler:           router,
		Addr:              apiAdress,
		TLSConfig:         nil,
		ReadTimeout:       time.Duration(readTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(readHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(writeTimeout) * time.Second,
		IdleTimeout:       time.Duration(idleTimeout) * time.Second,
		MaxHeaderBytes:    1 << 20,
		TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		ErrorLog:          log.New(os.Stderr, "SYSTEM ERROR:\t", log.LstdFlags),
		BaseContext: func(listener net.Listener) context.Context {
			return context.Background()
		},
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return ctx
		},
	}
}
