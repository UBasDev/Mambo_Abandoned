package models

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

	constants "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Constants"
	enums "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Enums"
	userControllers "github.com/UBasDev/Mambo/_src/MamboCoreService/Presentation/MamboCoreService.API/Controllers/UserController"
	helpers "github.com/UBasDev/Mambo/_src/_helpers"
	middlewares "github.com/UBasDev/Mambo/_src/_helpers/Middlewares"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type IApplicationBuild interface {
	ConfigureEnvironment(environment *enums.Environment)
	ConfigureApplicationSettingsFromConfigFile(environment enums.Environment, applicationConfig *ApplicationConfigModel) error
	InitZapLogger(isDevelopment bool, logFilePath string, debugLevel zapcore.Level, initialFields map[string]any, zapLogger **zap.Logger) error
	InitPostgreSqlDatabaseContext(postgreSqlDbUrl string, environment enums.Environment, postgreSqlDbContext **gorm.DB) error
	CreateRoutes(postgreSqlDbContext *gorm.DB, logger *zap.Logger, containerRouter **http.ServeMux, environment enums.Environment)
	CreateServer(containerServer *http.Server, router http.Handler, apiAdress string, readTimeout uint16, readHeaderTimeout uint16, writeTimeout uint16, idleTimeout uint16)
}
type ApplicationBuildModel struct {
}

func StartApplicationBuild() IApplicationBuild {
	return ApplicationBuildModel{}
}
func (applicationBuild ApplicationBuildModel) ConfigureEnvironment(environment *enums.Environment) {
	environmentFromOs := os.Getenv("GOLANG_ENVIRONMENT")
	if environmentFromOs == "" {
		*environment = enums.DevelopmentEnvironment
	} else {
		environment.Set(environmentFromOs)
	}
}
func (applicationBuild ApplicationBuildModel) ConfigureApplicationSettingsFromConfigFile(environment enums.Environment, applicationConfig *ApplicationConfigModel) error {
	var wg sync.WaitGroup
	wg.Add(1)
	applicationConfigChannel := make(chan ApplicationConfigModel)
	errorChannel := make(chan error)
	go helpers.ConfigFileRead[ApplicationConfigModel](environment, applicationConfigChannel, errorChannel, &wg)
	*applicationConfig = <-applicationConfigChannel
	errorValue := <-errorChannel
	wg.Wait()
	if errorValue != nil {
		return errorValue
	}
	return nil
}
func (applicationBuild ApplicationBuildModel) InitZapLogger(isDevelopment bool, logFilePath string, debugLevel zapcore.Level, initialFields map[string]any, zapLogger **zap.Logger) error {
	// Configure Zap logger options
	var config zap.Config
	if isDevelopment {
		config = zap.NewDevelopmentConfig()
		config.Encoding = "console"
	} else {
		config = zap.NewProductionConfig()
		config.Encoding = "json"
	}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.NameKey = "name"
	config.EncoderConfig.StacktraceKey = "stacktrace"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.LineEnding = "\n"
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Use ISO8601 time format
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.DisableCaller = false
	// Customize logging level for development
	config.Level.SetLevel(debugLevel)
	config.Development = isDevelopment
	config.OutputPaths = []string{"stdout", logFilePath}
	config.ErrorOutputPaths = []string{"stderr"}
	config.DisableStacktrace = false
	config.InitialFields = initialFields
	// Build the logger
	logger, err := config.Build()
	if err != nil {
		return err
	}
	defer logger.Sync() // Flushes buffer, if any
	*zapLogger = logger
	return nil
}

func (applicationBuild ApplicationBuildModel) InitPostgreSqlDatabaseContext(postgreSqlDbUrl string, environment enums.Environment, postgreSqlDbContext **gorm.DB) error {
	isParameterizedQueriesActive := true
	logLevel := logger.Error
	isLogColorful := false
	if environment != enums.ProductionEnvironment { //So if it is NOT production environment
		isParameterizedQueriesActive = false
		logLevel = logger.Info
		isLogColorful = true
	}
	newDatabaseLogger := logger.New(
		log.New(os.Stdout, "SQL:\t", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second * time.Duration(30),
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      isParameterizedQueriesActive,
			Colorful:                  isLogColorful,
		},
	)
	connectedDb, err := gorm.Open(postgres.Open(postgreSqlDbUrl), &gorm.Config{
		Logger:                 newDatabaseLogger,
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",   // table name prefix, table for `User` would be `t_users`
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   true,
		},
		AllowGlobalUpdate:    false,
		FullSaveAssociations: true,
		DryRun:               false,
		NowFunc:              time.Now().UTC,
		DisableAutomaticPing: false,
	})
	if err != nil {
		return err
	}
	//dB.AutoMigrate(&entities.RoleEntity{}, &entities.ScreenEntity{}, &entities.CompanyEntity{})
	//dB.AutoMigrate(&entities.UserEntity{}, &entities.ProfileEntity{})
	*postgreSqlDbContext = connectedDb
	return nil
}
func (applicationBuild ApplicationBuildModel) CreateRoutes(postgreSqlDbContext *gorm.DB, logger *zap.Logger, containerRouter **http.ServeMux, environment enums.Environment) {
	router := http.NewServeMux()
	createSingleUserController := userControllers.BuildSingleUserController(postgreSqlDbContext)
	router.Handle(fmt.Sprintf("/%s/%s/create-single-user", constants.ApiEndpointPrefix, constants.ApiVersion), middlewares.TraceIdMiddleware(createSingleUserController, logger, environment))
	*containerRouter = router
}
func (applicationBuild ApplicationBuildModel) CreateServer(containerServer *http.Server, router http.Handler, apiAdress string, readTimeout uint16, readHeaderTimeout uint16, writeTimeout uint16, idleTimeout uint16) {
	*containerServer = http.Server{
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
