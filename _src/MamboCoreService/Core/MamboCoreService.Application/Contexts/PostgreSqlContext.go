package contexts

import (
	"log"
	"os"
	"time"

	enums "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Enums"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitPostgreSqlDatabaseContext(postgreSqlDbUrl string, environment enums.Environment) (*gorm.DB, error) {
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
	dB, err := gorm.Open(postgres.Open(postgreSqlDbUrl), &gorm.Config{
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
		return nil, err
	}
	//dB.AutoMigrate(&entities.RoleEntity{}, &entities.ScreenEntity{}, &entities.CompanyEntity{})
	//dB.AutoMigrate(&entities.UserEntity{}, &entities.ProfileEntity{})

	return dB, nil
}
