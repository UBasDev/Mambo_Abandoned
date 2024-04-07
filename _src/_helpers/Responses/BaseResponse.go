package responses

import (
	"encoding/json"
	"net/http"
)

type IBaseResponse interface {
	SetErrorMessage(errorMessage string)
	SetTraceId(traceId string)
	SetServerTime(serverTime int64)
	SetStatusCode(statusCode uint16)
	SerializeForError() ([]byte, error)
}

type BaseResponse struct {
	issuccessful bool
	errorMessage string
	traceId      string
	serverTime   int64
	statusCode   uint16
}

func CreateNewBaseResponse() IBaseResponse {
	return &BaseResponse{
		issuccessful: true,
		statusCode:   http.StatusOK,
	}
}
func (baseResponse *BaseResponse) SetErrorMessage(errorMessage string) {
	baseResponse.errorMessage = errorMessage
	baseResponse.issuccessful = false
}
func (baseResponse *BaseResponse) SetTraceId(traceId string) {
	baseResponse.traceId = traceId
}
func (baseResponse *BaseResponse) SetServerTime(serverTime int64) {
	baseResponse.serverTime = serverTime
}
func (baseResponse *BaseResponse) SetStatusCode(statusCode uint16) {
	baseResponse.statusCode = statusCode
}
func (baseResponse *BaseResponse) SerializeForError() ([]byte, error) {
	serializedBaseResponse, err := json.Marshal(struct {
		Issuccessful bool
		ErrorMessage string
		TraceId      string
		ServerTime   int64
		StatusCode   uint16
	}{
		Issuccessful: false,
		ErrorMessage: baseResponse.errorMessage,
		TraceId:      baseResponse.traceId,
		ServerTime:   baseResponse.serverTime,
		StatusCode:   baseResponse.statusCode,
	})
	return serializedBaseResponse, err
}
