package main

import (
	"fmt"
	"log"

	enums "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Enums"
	models "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Models"
	"go.uber.org/zap/zapcore"
)

func main() {
	var environment enums.Environment
	var applicationConfig models.ApplicationConfigModel
	applicationBuilder := models.StartApplicationBuild()
	applicationBuilder.ConfigureEnvironment(&environment)
	err := applicationBuilder.ConfigureApplicationSettingsFromConfigFile(environment, &applicationConfig)
	if err != nil {
		log.Printf("failed to read configuration file. Error: %+v", err)
		panic(fmt.Sprintf("failed to read configuration file. Error: %+v", err))
	}
	applicationContainer := models.CreateNewApplicationContainerModel(
		fmt.Sprintf("%s:%s", applicationConfig.AppSettings.Server.Host, applicationConfig.AppSettings.Server.Port),
		applicationConfig.AppSettings.Server.ReadTimeout,
		applicationConfig.AppSettings.Server.ReadHeaderTimeout,
		applicationConfig.AppSettings.Server.WriteTimeout,
		applicationConfig.AppSettings.Server.IdleTimeout,
		environment,
	)
	err = applicationBuilder.InitZapLogger(true, applicationConfig.AppSettings.LogFilePath, zapcore.DebugLevel, map[string]any{}, &applicationContainer.ZapLogger)
	if err != nil {
		log.Printf("failed to initialize logger. Error: %+v", err)
		panic(fmt.Sprintf("failed to initialize logger. Error: %+v", err))
	}
	err = applicationBuilder.InitPostgreSqlDatabaseContext(applicationConfig.AppSettings.DatabaseConnectionStrings.PostgreSqlDbUrl, environment, &applicationContainer.PostgreSqlDbContext)
	if err != nil {
		applicationContainer.ZapLogger.Error(fmt.Sprintf("PostgreSql database connection couldnt established. Error: %+v", err))
		panic(fmt.Sprintf("PostgreSql database connection couldnt established. Error: %+v", err))
	}
	applicationBuilder.CreateRoutes(applicationContainer.PostgreSqlDbContext, applicationContainer.ZapLogger, &applicationContainer.Router, environment)
	applicationBuilder.CreateServer(&applicationContainer.Server, applicationContainer.Router, applicationContainer.ApiAdress, applicationContainer.ReadTimout, applicationContainer.ReadHeaderTimeout, applicationContainer.WriteTimeout, applicationContainer.IdleTimeout)
	err = applicationContainer.Server.ListenAndServe()
	if err != nil {
		applicationContainer.ZapLogger.Error(fmt.Sprintf("Server couldnt run. Error: %+v", err))
		panic(err)
	}
}
