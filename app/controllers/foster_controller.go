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

// Create a Foster
func CreateFosterHandler(respWr http.ResponseWriter, req *http.Request) {
	createFosterInput := types.CreateFosterInputDTO{}
	// Access restrictions
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Get Failed")
		return
	}
	if !access_utils.IsCrtUserAdmin(currentUser, respWr) {
		return
	}

	err := utils.ExtractReqBody(req, &createFosterInput)
	if err != nil {
		fmt.Println("[controller]", err.Error())
		utils.SendErrorResponse(respWr, 400, "Create Failed")
		return
	}

	createResult, createError := services.CreateFoster(createFosterInput)
	if createError == nil {
		utils.SendSuccessResponse(respWr, createResult)

	} else {
		utils.SendErrorResponse(respWr, 500, "Create Failed")
	}
}

// Get a Foster by ID
func GetFosterByIdHandler(respWr http.ResponseWriter, req *http.Request) {
	// Access restrictions
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Get Failed")
		return
	}
	if !access_utils.IsCrtUserAdmin(currentUser, respWr) {
		return
	}

	fosterId := req.PathValue("foster_id")
	idAsInt, err := strconv.Atoi(fosterId)
	if (idAsInt <= 0) || (err != nil) {
		fmt.Println("[controller] Invalid ID.", err.Error())
		utils.SendErrorResponse(respWr, 404, "Not Found")
		return
	}

	getByIdResult, getByIdError := services.GetFosterById(idAsInt)
	if getByIdError == nil {
		utils.SendSuccessResponse(respWr, getByIdResult)

	} else if utils.IsErrorOfType(getByIdError, utils.ErrorType_FosterDoesNotExist) {
		utils.SendErrorResponse(respWr, 404, "Not Found")

	} else {
		utils.SendErrorResponse(respWr, 500, "Get Failed")
	}
}

// Get all Fosters (+ Query)
func GetAllFostersHandler(respWr http.ResponseWriter, req *http.Request) {
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
	queryInfo, queryError := query_utils.ParseQueryParams[models.Foster](query)
	if utils.IsErrorOfType(queryError, utils.ErrorType_QueryError) {
		fmt.Println("[controller]", queryError.Error())
		utils.SendErrorResponse(respWr, 400, fmt.Sprintf("Bad Input. %s", queryError.Error()))
		return
	}

	getAllFostersResult, getAllFostersError := services.GetAllFosters(queryInfo)
	if getAllFostersError == nil {
		utils.SendSuccessResponse(respWr, getAllFostersResult)

	} else if utils.IsErrorOfType(getAllFostersError, utils.ErrorType_QueryError) {
		utils.SendErrorResponse(respWr, 400, fmt.Sprintf("Bad Input. %s", getAllFostersError.Error()))

	} else {
		utils.SendErrorResponse(respWr, 500, "Get Failed")
	}
}

// Update a Foster by ID
func UpdateFosterHandler(respWr http.ResponseWriter, req *http.Request) {
	// Access restrictions
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Update Failed")
		return
	}
	if !access_utils.IsCrtUserAdmin(currentUser, respWr) {
		return
	}

	fosterId := req.PathValue("foster_id")
	idAsInt, err := strconv.Atoi(fosterId)
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

	updateResult, updateError := services.UpdateFoster(idAsInt, updateInput)
	if updateError == nil {
		utils.SendSuccessResponse(respWr, updateResult)

	} else if utils.IsErrorOfType(updateError, utils.ErrorType_FosterDoesNotExist) {
		utils.SendErrorResponse(respWr, 404, "Not Found")

	} else {
		utils.SendErrorResponse(respWr, 500, "Update Failed")
	}
}

// Delete a Foster by ID
func DeleteFosterHandler(respWr http.ResponseWriter, req *http.Request) {
	// Access restrictions
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Delete Failed")
		return
	}
	if !access_utils.IsCrtUserAdmin(currentUser, respWr) {
		return
	}

	fosterId := req.PathValue("foster_id")
	idAsInt, err := strconv.Atoi(fosterId)
	if (idAsInt <= 0) || (err != nil) {
		fmt.Println("[controller] Invalid ID.", err.Error())
		utils.SendErrorResponse(respWr, 404, "Not Found")
		return
	}

	deleteResult, deleteError := services.DeleteFoster(idAsInt)
	if deleteError == nil {
		utils.SendSuccessResponse(respWr, deleteResult)

	} else if utils.IsErrorOfType(deleteError, utils.ErrorType_FosterDoesNotExist) {
		utils.SendErrorResponse(respWr, 404, "Not Found")

	} else {
		utils.SendErrorResponse(respWr, 500, "Delete Failed")
	}
}
