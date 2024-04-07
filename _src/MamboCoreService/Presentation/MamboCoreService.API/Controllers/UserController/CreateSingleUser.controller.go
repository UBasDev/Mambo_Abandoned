package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	models "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/Models"
	requestmodels "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Application/RequestModels"
	entities "github.com/UBasDev/Mambo/_src/MamboCoreService/Core/MamboCoreService.Domain/Entities"
	responses "github.com/UBasDev/Mambo/_src/_helpers/Responses"
	"gorm.io/gorm"
)

type ICreateSingleUserController interface {
	ServeHTTP(rw http.ResponseWriter, rq *http.Request)
}

type CreateSingleUserControllerModel struct {
	postgreSqlDbContext *gorm.DB
}

func BuildSingleUserController(postgreSqlDbContext *gorm.DB) ICreateSingleUserController {
	return &CreateSingleUserControllerModel{
		postgreSqlDbContext,
	}
}
func (createSingleUserControllerModel CreateSingleUserControllerModel) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	response := responses.CreateNewBaseResponse()
	traceId, ok := rq.Context().Value(models.TraceIdKey{}).(string)
	if !ok {
		log.Println("TraceId couldn't be generated")
		responses.GenerateErrorResponse(rw, response, "", "TraceId couldn't be generated", http.StatusInternalServerError)
		return
	}
	if rq.Method != http.MethodPost {
		log.Printf("Wrong HTTP request method: %s", rq.Method)
		responses.GenerateErrorResponse(rw, response, traceId, fmt.Sprintf("Wrong HTTP request method: %s", rq.Method), http.StatusBadRequest)
		return
	}
	var requestBody requestmodels.CreateSingleUserRequestModel

	if err := json.NewDecoder(rq.Body).Decode(&requestBody); err != nil {
		log.Printf("Check your request body, unable to decode: %s", err)
		responses.GenerateErrorResponse(rw, response, traceId, "Request body is invalid", http.StatusBadRequest)
		return
	}
	if err := requestBody.Validate(); err != nil {
		log.Printf("Request body is invalid: %s", err)
		responses.GenerateErrorResponse(rw, response, traceId, fmt.Sprintf("Request body is invalid: %s", err), http.StatusBadRequest)
		return
	}

	userToCreate := entities.BuildNewUserEntity(requestBody.Username, requestBody.Email)
	resultsFromUserCreate := createSingleUserControllerModel.postgreSqlDbContext.Create(userToCreate)
	if resultsFromUserCreate.Error != nil {
		log.Printf("Unable to write this user object to the database: %s", resultsFromUserCreate.Error)
		responses.GenerateErrorResponse(rw, response, traceId, "Unable to write this user object to the database", http.StatusBadRequest)
		return
	}
}
