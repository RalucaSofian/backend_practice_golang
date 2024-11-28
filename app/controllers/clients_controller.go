package controllers

import (
	"app/models"
	"app/services"
	"app/types"
	"app/utils"
	"app/utils/access_utils"
	"app/utils/middlewares"
	"app/utils/query_utils"
	"fmt"
	"net/http"
	"strconv"
)

// Create a Client
func CreateClientHandler(respWr http.ResponseWriter, req *http.Request) {
	createClientInput := types.CreateClientInputDTO{}

	// Access restrictions
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Create Failed")
		return
	}
	if !access_utils.IsCrtUserAdmin(currentUser, respWr) {
		return
	}

	err := utils.ExtractReqBody(req, &createClientInput)
	if err != nil {
		fmt.Println("[controller]", err.Error())
		utils.SendErrorResponse(respWr, 400, "Create Failed")
		return
	}

	createResult, createError := services.CreateClient(createClientInput)
	if createError == nil {
		utils.SendSuccessResponse(respWr, createResult)

	} else {
		utils.SendErrorResponse(respWr, 500, "Create Failed")
	}
}

// Get a Client by ID
func GetClientByIdHandler(respWr http.ResponseWriter, req *http.Request) {
	// Access restrictions
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Get Failed")
		return
	}
	if !access_utils.IsCrtUserAdmin(currentUser, respWr) {
		return
	}

	clientId := req.PathValue("client_id")
	idAsInt, err := strconv.Atoi(clientId)
	if (idAsInt <= 0) || (err != nil) {
		fmt.Println("[controller] Invalid ID.", err.Error())
		utils.SendErrorResponse(respWr, 404, "Not Found")
		return
	}

	getByIdResult, getByIdError := services.GetClientById(idAsInt)
	if getByIdError == nil {
		utils.SendSuccessResponse(respWr, getByIdResult)

	} else if utils.IsErrorOfType(getByIdError, utils.ErrorType_ClientDoesNotExist) {
		utils.SendErrorResponse(respWr, 404, "Not Found")

	} else {
		utils.SendErrorResponse(respWr, 500, "Get Failed")
	}
}

// Get all Clients (+ Query)
func GetAllClientsHandler(respWr http.ResponseWriter, req *http.Request) {
	// Access restrictions
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Get Failed")
		return
	}
	if !access_utils.IsCrtUserAdmin(currentUser, respWr) {
		return
	}

	query := req.URL.Query()
	queryInfo, queryError := query_utils.ParseQueryParams[models.Client](query)
	if utils.IsErrorOfType(queryError, utils.ErrorType_QueryError) {
		fmt.Println("[controller]", queryError.Error())
		utils.SendErrorResponse(respWr, 400, fmt.Sprintf("Bad Input. %s", queryError.Error()))
		return
	}

	getAllClientsResult, getAllClientsError := services.GetAllClients(queryInfo)
	if getAllClientsError == nil {
		utils.SendSuccessResponse(respWr, getAllClientsResult)

	} else if utils.IsErrorOfType(getAllClientsError, utils.ErrorType_QueryError) {
		utils.SendErrorResponse(respWr, 400, fmt.Sprintf("Bad Input. %s", getAllClientsError.Error()))

	} else {
		utils.SendErrorResponse(respWr, 500, "Get Failed")
	}
}

// Update a Client by ID
func UpdateClientHandler(respWr http.ResponseWriter, req *http.Request) {
	// Access restrictions
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Update Failed")
		return
	}
	if !access_utils.IsCrtUserAdmin(currentUser, respWr) {
		return
	}

	clientId := req.PathValue("client_id")
	idAsInt, err := strconv.Atoi(clientId)
	if (idAsInt <= 0) || (err != nil) {
		fmt.Println("[controller] Invalid ID.", err.Error())
		utils.SendErrorResponse(respWr, 404, "Not Found")
		return
	}

	updateInput := make(map[string]interface{})
	err = utils.ExtractReqBody(req, &updateInput)
	if err != nil {
		fmt.Println("[controller]", err.Error())
		utils.SendErrorResponse(respWr, 400, "Update Failed")
		return
	}

	updateResult, updateError := services.UpdateClient(idAsInt, updateInput)
	if updateError == nil {
		utils.SendSuccessResponse(respWr, updateResult)

	} else if utils.IsErrorOfType(updateError, utils.ErrorType_ClientDoesNotExist) {
		utils.SendErrorResponse(respWr, 404, "Not Found")

	} else {
		utils.SendErrorResponse(respWr, 500, "Update Failed")
	}
}

// Delete a Client by ID
func DeleteClientHandler(respWr http.ResponseWriter, req *http.Request) {
	// Access restrictions
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Delete Failed")
		return
	}
	if !access_utils.IsCrtUserAdmin(currentUser, respWr) {
		return
	}

	clientId := req.PathValue("client_id")
	idAsInt, err := strconv.Atoi(clientId)
	if (idAsInt <= 0) || (err != nil) {
		fmt.Println("[controller] Invalid ID.", err.Error())
		utils.SendErrorResponse(respWr, 404, "Not Found")
		return
	}

	deleteResult, deleteError := services.DeleteClient(idAsInt)
	if deleteError == nil {
		utils.SendSuccessResponse(respWr, deleteResult)

	} else if utils.IsErrorOfType(deleteError, utils.ErrorType_ClientDoesNotExist) {
		utils.SendErrorResponse(respWr, 404, "Not Found")

	} else {
		utils.SendErrorResponse(respWr, 500, "Delete Failed")
	}
}
