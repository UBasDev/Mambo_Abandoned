package responses

import (
	"log"
	"net/http"
	"time"
)

func GenerateErrorResponse(rw http.ResponseWriter, response IBaseResponse, traceId string, errorMessage string, statusCode uint16) {
	response.SetErrorMessage(errorMessage)
	response.SetStatusCode(statusCode)
	response.SetServerTime(time.Now().Unix())
	response.SetTraceId(traceId)
	serializedResponse, err := response.SerializeForError()
	if err != nil {
		log.Printf("Something went wrong while serializing response: %s", err)
		http.Error(rw, "We couldn't send the response correctly", http.StatusInternalServerError)
		return
	}
	_, err = rw.Write(serializedResponse)
	if err != nil {
		log.Printf("Something went wrong while writing response: %s", err)
		http.Error(rw, "We couldn't send the response correctly", http.StatusInternalServerError)
		return
	}
}
