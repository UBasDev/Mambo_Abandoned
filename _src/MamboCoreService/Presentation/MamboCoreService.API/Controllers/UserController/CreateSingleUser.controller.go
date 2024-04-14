package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	requestmodels "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/RequestModels"
	entities "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Domain/Entities"
	middlewares "github.com/UBasDev/Mambo/_src/_helpers/Middlewares"
	responses "github.com/UBasDev/Mambo/_src/_helpers/Responses"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ICreateSingleUserController interface {
	ServeHTTP(rw http.ResponseWriter, rq *http.Request)
}

type CreateSingleUserControllerModel struct {
	postgreSqlDbContext *gorm.DB
	zipkinClient        *zipkinhttp.Client
}

func BuildSingleUserController(postgreSqlDbContext *gorm.DB, zipkinClient *zipkinhttp.Client) ICreateSingleUserController {
	return &CreateSingleUserControllerModel{
		postgreSqlDbContext,
		zipkinClient,
	}
}
func (createSingleUserControllerModel CreateSingleUserControllerModel) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	response := responses.CreateNewBaseResponse()
	logContext, ok := rq.Context().Value(middlewares.ZapLogger{}).(*zap.Logger)
	if !ok {
		log.Println("Log context couldn't be generated")
		responses.GenerateErrorResponse(rw, response, "", "LogContext couldn't be generated", http.StatusInternalServerError)
		return
	}
	traceId, ok := rq.Context().Value(middlewares.TraceIdKey{}).(string)
	if !ok {
		logContext.Error("TraceId couldn't be generated")
		responses.GenerateErrorResponse(rw, response, "", "TraceId couldn't be generated", http.StatusInternalServerError)
		return
	}
	//ZIPKIN
	newRequest, err := http.NewRequest(http.MethodGet, "https://pokeapi.co/api/v2/pokemon/ditto", nil)
	if err != nil {
		logContext.Error(fmt.Sprintf("Error: %+V", err))
		return
	}
	span := zipkin.SpanFromContext(rq.Context())
	ctx := zipkin.NewContext(newRequest.Context(), span)

	newRequest = newRequest.WithContext(ctx)
	resp, err := createSingleUserControllerModel.zipkinClient.Do(newRequest)
	if err != nil {
		logContext.Error(fmt.Sprintf("Error: %+V", err))
		return
	}
	logContext.Debug(fmt.Sprintf("Response: %+V", resp))
	childSpan := zipkin.SpanFromContext(ctx)
	childSpan.Tag("child1", "childValue1")
	//ZIPKIN

	if rq.Method != http.MethodPost {
		logContext.Error(fmt.Sprintf("Wrong HTTP request method: %s", rq.Method))
		responses.GenerateErrorResponse(rw, response, traceId, fmt.Sprintf("Wrong HTTP request method: %s", rq.Method), http.StatusBadRequest)
		return
	}
	var requestBody requestmodels.CreateSingleUserRequestModel

	if err := json.NewDecoder(rq.Body).Decode(&requestBody); err != nil {
		logContext.Error(fmt.Sprintf("Check your request body, unable to decode: %s", err))
		responses.GenerateErrorResponse(rw, response, traceId, "Request body is invalid", http.StatusBadRequest)
		return
	}
	if err := requestBody.Validate(); err != nil {
		logContext.Error(fmt.Sprintf("Request body is invalid: %s", err))
		responses.GenerateErrorResponse(rw, response, traceId, fmt.Sprintf("Request body is invalid: %s", err), http.StatusBadRequest)
		return
	}

	userToCreate := entities.BuildNewUserEntity(requestBody.Username, requestBody.Email)
	resultsFromUserCreate := createSingleUserControllerModel.postgreSqlDbContext.Create(userToCreate)
	if resultsFromUserCreate.Error != nil {
		logContext.Error(fmt.Sprintf("Unable to write this user object to the database: %s", resultsFromUserCreate.Error))
		responses.GenerateErrorResponse(rw, response, traceId, "Unable to write this user object to the database", http.StatusBadRequest)
		return
	}
}
