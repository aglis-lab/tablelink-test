package response

import (
	"context"
	"encoding/json"
	"net/http"
	"tablelink/src/middleware"
)

const (
	BadRequest     = "bad request"
	InternalServer = "internal server error"
)

type Error struct {
	Code     string `json:"code"`
	Title    string `json:"message_title"`
	Message  string `json:"message"`
	Severity string `json:"message_severity"`
}

type Meta struct {
	RequestId string `json:"request_id"`
}

type Response struct {
	Data     interface{} `json:"data"`
	Error    *Error      `json:"error"`
	Metadata Meta        `json:"metadata"`
}

func createErrorResponse(err, reqId string) Response {
	return Response{
		Error: &Error{
			Message:  err,
			Severity: "error",
		},
		Metadata: Meta{
			RequestId: reqId,
		},
	}
}

func createSuccessResponse(data interface{}, reqId string) Response {
	return Response{
		Data: data,
		Metadata: Meta{
			RequestId: reqId,
		},
	}
}

func JSONResponse(w http.ResponseWriter, data Response, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func JSONBadRequestResponse(ctx context.Context, w http.ResponseWriter) {
	JSONResponse(w, createErrorResponse(BadRequest, middleware.GetRequestID(ctx)),
		http.StatusBadRequest)
}

func JSONInternalServerError(ctx context.Context, w http.ResponseWriter) {
	JSONResponse(w, createErrorResponse(InternalServer, middleware.GetRequestID(ctx)),
		http.StatusInternalServerError)
}

func JSONSuccessResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	JSONResponse(w, createSuccessResponse(data, middleware.GetRequestID(ctx)), http.StatusOK)
}
