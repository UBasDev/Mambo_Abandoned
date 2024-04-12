package routes

import (
	"fmt"
	"net/http"

	enums "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Enums"
	userControllers "github.com/UBasDev/Mambo/_src/MamboCoreService/Presentation/MamboCoreService.API/Controllers/UserController"
	middlewares "github.com/UBasDev/Mambo/_src/_helpers/Middlewares"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	apiEndpointPrefix string = "api"
	apiVersion        string = "v1"
)

func CreateRoutes(postgreSqlDbContext *gorm.DB, logger *zap.Logger, environment enums.Environment) *http.ServeMux {
	router := http.NewServeMux()
	createSingleUserController := userControllers.BuildSingleUserController(postgreSqlDbContext)
	router.Handle(fmt.Sprintf("/%s/%s/create-single-user", apiEndpointPrefix, apiVersion), middlewares.TraceIdMiddleware(createSingleUserController, logger, environment))
	return router
}
