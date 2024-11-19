package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type responseBody struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func sendResponse(respWr http.ResponseWriter, statusCode int, message string, data any) {
	response := responseBody{Message: message, Data: data}
	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println("[utils]", err.Error())
	}

	respWr.WriteHeader(statusCode)
	respWr.Write(bytes)
}

func SendSuccessResponse(respWr http.ResponseWriter, data any) {
	sendResponse(respWr, 200, "Success", data)
}

func SendErrorResponse(respWr http.ResponseWriter, errCode int, message string) {
	sendResponse(respWr, errCode, message, nil)
}
