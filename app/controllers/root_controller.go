package controllers

import (
	"app/utils"
	"fmt"
	"net/http"
)

// Handle the HTTP Requests to root endpoint
func RootHandler(respWr http.ResponseWriter, req *http.Request) {
	utils.SendSuccessResponse(respWr, "Server Up")
	fmt.Println("[controller] Root endpoint request: ", req.Method, req.RequestURI)
}
