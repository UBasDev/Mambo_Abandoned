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

	constants "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Constants"
	enums "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Enums"
	userControllers "github.com/UBasDev/Mambo/_src/MamboCoreService/Presentation/MamboCoreService.API/Controllers/UserController"
	helpers "github.com/UBasDev/Mambo/_src/_helpers"
	middlewares "github.com/UBasDev/Mambo/_src/_helpers/Middlewares"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type IApplicationBuild interface {
	ConfigureEnvironment(environment *enums.Environment)
	ConfigureApplicationSettingsFromConfigFile(environment enums.Environment, applicationConfig *applicationConfigModel) error
	InitZapLogger(isDevelopment bool, logFilePath string, logLevel zapcore.Level, initialFields map[string]any, zapLogger **zap.Logger) error
	InitZapLoggerWithLoki(isDevelopment bool, logFilePath string, logLevel zapcore.Level, initialFields map[string]any, zapLogger **zap.Logger) error
	InitPostgreSqlDatabaseContext(postgreSqlDbUrl string, environment enums.Environment, postgreSqlDbContext **gorm.DB) error
	CreateRoutes(postgreSqlDbContext *gorm.DB, logger *zap.Logger, containerRouter **http.ServeMux)
	CreateServer(containerServer *http.Server, router http.Handler, apiAdress string, readTimeout uint16, readHeaderTimeout uint16, writeTimeout uint16, idleTimeout uint16)
}
type ApplicationBuildModel struct {
}

func startApplicationBuild() IApplicationBuild {
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
func (applicationBuild ApplicationBuildModel) ConfigureApplicationSettingsFromConfigFile(environment enums.Environment, applicationConfig *applicationConfigModel) error {
	var wg sync.WaitGroup
	wg.Add(1)
	applicationConfigChannel := make(chan applicationConfigModel)
	errorChannel := make(chan error)
	go helpers.ConfigFileRead[applicationConfigModel](environment, applicationConfigChannel, errorChannel, &wg)
	*applicationConfig = <-applicationConfigChannel
	errorValue := <-errorChannel
	wg.Wait()
	if errorValue != nil {
		return errorValue
	}
	return nil
}
func (applicationBuild ApplicationBuildModel) InitZapLogger(isDevelopment bool, logFilePath string, logLevel zapcore.Level, initialFields map[string]any, zapLogger **zap.Logger) error {

	// Configure Zap logger options
	var config zap.Config
	if isDevelopment {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}
	config.Encoding = "json"
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
	config.Level.SetLevel(logLevel)
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
func (applicationBuild ApplicationBuildModel) InitZapLoggerWithLoki(isDevelopment bool, logFilePath string, logLevel zapcore.Level, initialFields map[string]any, zapLogger **zap.Logger) error {

	// Configure Zap logger options
	var config zap.Config
	if isDevelopment {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}
	config.Encoding = "json"
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
	config.Level.SetLevel(logLevel)
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
	//connectedDb.AutoMigrate(&entities.RoleEntity{}, &entities.ScreenEntity{}, &entities.CompanyEntity{})
	//connectedDb.AutoMigrate(&entities.UserEntity{}, &entities.ProfileEntity{})
	*postgreSqlDbContext = connectedDb
	return nil
}
func (applicationBuild ApplicationBuildModel) CreateRoutes(postgreSqlDbContext *gorm.DB, logger *zap.Logger, containerRouter **http.ServeMux) {
	router := http.NewServeMux()
	tracer, err := NewTracer("mambo-core")
	if err != nil {
		logger.Error(fmt.Sprintf("Tracer couldn't be created: %+v", err))
	}
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer, zipkinhttp.TagResponseSize(true),
	)
	zipkinClient, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		logger.Error("unable to create client: %+v\n", err)
	}
	createSingleUserController := userControllers.BuildSingleUserController(postgreSqlDbContext, zipkinClient)
	router.Handle(fmt.Sprintf("/%s/%s/create-single-user", constants.ApiEndpointPrefix, constants.ApiVersion), serverMiddleware(middlewares.TraceIdMiddleware(createSingleUserController, logger)))
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
func NewTracer(serviceName string) (*zipkin.Tracer, error) {

	reporter := httpreporter.NewReporter("http://localhost:9411/api/v2/spans")

	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: serviceName, IPv4: getOutboundIP(), Port: 8080}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00)
	// of traces.

	sampler, err := zipkin.NewBoundarySampler(1, 10000000)
	if err != nil {
		return nil, err
	}
	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
		zipkin.WithTags(map[string]string{
			"environment": "development",
		}),
		zipkin.WithTraceID128Bit(true),
	)
	if err != nil {
		return nil, err
	}

	return tracer, err
}

// Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
