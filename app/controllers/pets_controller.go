package controllers

import (
	"app/models"
	"app/services"
	"app/types"
	"app/utils"
	"app/utils/middlewares"
	"app/utils/query_utils"
	"fmt"
	"net/http"
	"strconv"
)

// Create a Pet
func CreatePetHandler(respWr http.ResponseWriter, req *http.Request) {
	createPetInput := types.CreatePetInputDTO{}

	// Access restrictions
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Create Failed")
		return
	}
	fmt.Println("[controller] Current User is:", currentUser.Role)
	if currentUser.Role != models.UserRole_Admin {
		fmt.Println("[controller] Unauthorized request")
		utils.SendErrorResponse(respWr, 401, "Unauthorized")
		return
	}

	err := utils.ExtractReqBody(req, &createPetInput)
	if err != nil {
		fmt.Println("[controller]", err.Error())
		utils.SendErrorResponse(respWr, 400, "Create Failed")
		return
	}

	createResult, createError := services.CreatePet(createPetInput)
	if createError == nil {
		utils.SendSuccessResponse(respWr, createResult)

	} else {
		utils.SendErrorResponse(respWr, 500, "Create Failed")
	}
}

// Get a Pet by ID
func GetPetByIdHandler(respWr http.ResponseWriter, req *http.Request) {
	petId := req.PathValue("pet_id")
	idAsInt, err := strconv.Atoi(petId)
	if (idAsInt <= 0) || (err != nil) {
		fmt.Println("[controller] Invalid ID.", err.Error())
		utils.SendErrorResponse(respWr, 404, "Not Found")
		return
	}

	getByIdResult, getByIdError := services.GetPetById(idAsInt)
	if getByIdError == nil {
		utils.SendSuccessResponse(respWr, getByIdResult)

	} else if utils.IsErrorOfType(getByIdError, utils.ErrorType_PetDoesNotExist) {
		utils.SendErrorResponse(respWr, 404, "Not Found")

	} else {
		utils.SendErrorResponse(respWr, 500, "Get Failed")
	}
}

// Get all Pets (+ Query)
func GetAllPetsHandler(respWr http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	queryInfo, queryError := query_utils.ParseQueryParams[models.Pet](query)
	if utils.IsErrorOfType(queryError, utils.ErrorType_QueryError) {
		fmt.Println("[controller]", queryError.Error())
		utils.SendErrorResponse(respWr, 400, fmt.Sprintf("Bad Input. %s", queryError.Error()))
		return
	}

	getAllPetsResult, getAllPetsError := services.GetAllPets(queryInfo)
	if getAllPetsError == nil {
		utils.SendSuccessResponse(respWr, getAllPetsResult)

	} else if utils.IsErrorOfType(getAllPetsError, utils.ErrorType_QueryError) {
		utils.SendErrorResponse(respWr, 400, fmt.Sprintf("Bad Input. %s", getAllPetsError.Error()))

	} else {
		utils.SendErrorResponse(respWr, 500, "Get Failed")
	}
}

// Update a Pet by ID
func UpdatePetHandler(respWr http.ResponseWriter, req *http.Request) {
	// Access restrictions
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Update Failed")
		return
	}
	fmt.Println("[controller] Current User is:", currentUser.Role)
	if currentUser.Role != models.UserRole_Admin {
		fmt.Println("[controller] Unauthorized request")
		utils.SendErrorResponse(respWr, 401, "Unauthorized")
		return
	}

	petId := req.PathValue("pet_id")
	idAsInt, err := strconv.Atoi(petId)
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

	updateResult, updateError := services.UpdatePet(idAsInt, updateInput)
	if updateError == nil {
		utils.SendSuccessResponse(respWr, updateResult)

	} else if utils.IsErrorOfType(updateError, utils.ErrorType_PetDoesNotExist) {
		utils.SendErrorResponse(respWr, 404, "Not Found")

	} else {
		utils.SendErrorResponse(respWr, 500, "Update Failed")
	}
}

// Delete a Pet by ID
func DeletePetHandler(respWr http.ResponseWriter, req *http.Request) {
	// Access restrictions
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Update Failed")
		return
	}
	fmt.Println("[controller] Current User is:", currentUser.Role)
	if currentUser.Role != models.UserRole_Admin {
		fmt.Println("[controller] Unauthorized request")
		utils.SendErrorResponse(respWr, 401, "Unauthorized")
		return
	}

	petId := req.PathValue("pet_id")
	idAsInt, err := strconv.Atoi(petId)
	if (idAsInt <= 0) || (err != nil) {
		fmt.Println("[controller] Invalid ID.", err.Error())
		utils.SendErrorResponse(respWr, 404, "Not Found")
		return
	}

	deleteResult, deleteError := services.DeletePet(idAsInt)
	if deleteError == nil {
		utils.SendSuccessResponse(respWr, deleteResult)

	} else if utils.IsErrorOfType(deleteError, utils.ErrorType_PetDoesNotExist) {
		utils.SendErrorResponse(respWr, 404, "Not Found")

	} else {
		utils.SendErrorResponse(respWr, 500, "Delete Failed")
	}
}
