package routes

import (
	"fmt"
	"net/http"

	userControllers "github.com/UBasDev/Mambo/_src/MamboCoreService/Presentation/MamboCoreService.API/Controllers/UserController"
	middlewares "github.com/UBasDev/Mambo/_src/_helpers/Middlewares"
	"gorm.io/gorm"
)

const (
	apiEndpointPrefix string = "api"
	apiVersion        string = "v1"
)

func CreateRoutes(postgreSqlDbContext *gorm.DB) *http.ServeMux {
	router := http.NewServeMux()
	createSingleUserController := userControllers.BuildSingleUserController(postgreSqlDbContext)
	router.Handle(fmt.Sprintf("/%s/%s/create-single-user", apiEndpointPrefix, apiVersion), middlewares.TraceIdMiddleware(createSingleUserController))
	return router
}
