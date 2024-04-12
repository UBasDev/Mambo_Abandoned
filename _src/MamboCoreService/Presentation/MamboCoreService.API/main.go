package main

import (
	"context"
	"crypto/tls"

	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	contexts "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Contexts"
	enums "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Enums"
	models "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Models"
	routes "github.com/UBasDev/Mambo/_src/MamboCoreService/Presentation/MamboCoreService.API/Routes"
	helpers "github.com/UBasDev/Mambo/_src/_helpers"
	logproviders "github.com/UBasDev/Mambo/_src/_helpers/LogProviders"
	"go.uber.org/zap/zapcore"
)

func main() {
	var environment enums.Environment
	configureEnvironment(&environment)
	var wg sync.WaitGroup
	wg.Add(1)
	applicationConfigChannel := make(chan models.ApplicationConfig)
	go helpers.ConfigFileRead[models.ApplicationConfig](environment, applicationConfigChannel, &wg)
	applicationConfig := <-applicationConfigChannel
	wg.Wait()
	applicationContainer := models.ApplicationContainerModel{
		ApiAdress:         fmt.Sprintf("%s:%s", applicationConfig.AppSettings.Server.Host, applicationConfig.AppSettings.Server.Port),
		ReadTimout:        applicationConfig.AppSettings.Server.ReadTimeout,
		ReadHeaderTimeout: applicationConfig.AppSettings.Server.ReadHeaderTimeout,
		WriteTimeout:      applicationConfig.AppSettings.Server.WriteTimeout,
		IdleTimeout:       applicationConfig.AppSettings.Server.IdleTimeout,
		Environment:       environment,
	}
	logger := logproviders.InitZapLogger(true, applicationConfig.AppSettings.LogFilePath, zapcore.DebugLevel, map[string]any{})
	postgreSqlDbContext, err := contexts.InitPostgreSqlDatabaseContext(applicationConfig.AppSettings.DatabaseConnectionStrings.PostgreSqlDbUrl, environment)
	if err != nil {
		logger.Error(fmt.Sprintf("PostgreSql database connection couldnt established. Error: %+v", err))
		panic(err)
	}
	router := routes.CreateRoutes(postgreSqlDbContext, logger, environment)
	server := createServer(router, applicationContainer.ApiAdress, applicationContainer.ReadTimout, applicationContainer.ReadHeaderTimeout, applicationContainer.WriteTimeout, applicationContainer.IdleTimeout)
	err = server.ListenAndServe()
	if err != nil {
		logger.Error(fmt.Sprintf("Server couldnt run. Error: %+v", err))
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
func configureEnvironment(environment *enums.Environment) {
	environmentFromOs := os.Getenv("GOLANG_ENVIRONMENT")
	if environmentFromOs == "" {
		*environment = enums.ProductionEnvironment
	} else {
		environment.Set(environmentFromOs)
	}
}
