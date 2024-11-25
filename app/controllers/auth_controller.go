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

// Register User
func RegisterHandler(respWr http.ResponseWriter, req *http.Request) {
	registerInput := types.RegisterInputDTO{}

	err := utils.ExtractReqBody(req, &registerInput)
	if err != nil {
		fmt.Println("[controller]", err.Error())
		utils.SendErrorResponse(respWr, 400, "Register Failed")
		return
	}

	registerResult, registerError := services.RegisterUser(registerInput)
	if registerError == nil {
		utils.SendSuccessResponse(respWr, registerResult)

	} else if utils.IsErrorOfType(registerError, utils.ErrorType_UserAlreadyExists) {
		utils.SendErrorResponse(respWr, 400, "User Already Exists")

	} else {
		utils.SendErrorResponse(respWr, 500, "Register Failed")
	}
}

// Login User
func LoginHandler(respWr http.ResponseWriter, req *http.Request) {
	loginInput := types.LoginInputDTO{}

	err := utils.ExtractReqBody(req, &loginInput)
	if err != nil {
		fmt.Println("[controller]", err.Error())
		utils.SendErrorResponse(respWr, 400, "Login Failed")
		return
	}

	loginResult, loginError := services.LoginUser(loginInput)
	if loginError == nil {
		utils.SendSuccessResponse(respWr, loginResult)

	} else if utils.IsErrorOfType(loginError, utils.ErrorType_UserLoginFailed) {
		utils.SendErrorResponse(respWr, 401, "Unauthorized")

	} else {
		utils.SendErrorResponse(respWr, 500, "Login Failed")
	}
}

// Get an User by ID
func GetUserByIdHandler(respWr http.ResponseWriter, req *http.Request) {
	userId := req.PathValue("user_id")
	idAsInt, err := strconv.Atoi(userId)
	if (idAsInt <= 0) || (err != nil) {
		fmt.Println("[controller] Invalid ID.", err.Error())
		utils.SendErrorResponse(respWr, 404, "Not Found")
		return
	}

	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Get Failed")
		return
	}

	fmt.Println("[controller] Current User is:", currentUser.Role)
	if currentUser.Role != models.UserRole_Admin {
		if currentUser.Id != idAsInt {
			fmt.Println("[controller] Unauthorized request on Other User")
			utils.SendErrorResponse(respWr, 401, "Unauthorized")
			return
		}
	}

	getByIdResult, getByIdError := services.GetUserById(idAsInt)
	if getByIdError == nil {
		utils.SendSuccessResponse(respWr, getByIdResult)

	} else if utils.IsErrorOfType(getByIdError, utils.ErrorType_UserDoesNotExist) {
		utils.SendErrorResponse(respWr, 404, "Not Found")

	} else {
		utils.SendErrorResponse(respWr, 500, "Get Failed")
	}
}

// Get all Users (+ Query)
func GetAllUsersHandler(respWr http.ResponseWriter, req *http.Request) {
	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Get Failed")
		return
	}

	fmt.Println("[controller] Current User is:", currentUser.Role)
	if currentUser.Role != models.UserRole_Admin {
		fmt.Println("[controller] Unauthorized request on Other Users")
		utils.SendErrorResponse(respWr, 401, "Unauthorized")
		return
	}

	query := req.URL.Query()
	queryInfo, queryError := query_utils.ParseQueryParams[models.AuthUser](query)
	if utils.IsErrorOfType(queryError, utils.ErrorType_QueryError) {
		fmt.Println("[controller]", queryError.Error())
		utils.SendErrorResponse(respWr, 400, fmt.Sprintf("Bad Input. %s", queryError.Error()))
		return
	}

	getAllUsersResult, getAllUsersError := services.GetAllUsers(queryInfo)
	if getAllUsersError == nil {
		utils.SendSuccessResponse(respWr, getAllUsersResult)

	} else if utils.IsErrorOfType(getAllUsersError, utils.ErrorType_QueryError) {
		utils.SendErrorResponse(respWr, 400, fmt.Sprintf("Bad Input. %s", getAllUsersError.Error()))

	} else {
		utils.SendErrorResponse(respWr, 500, "Get Failed")
	}
}

// Update an User by ID
func UpdateUserHandler(respWr http.ResponseWriter, req *http.Request) {
	userId := req.PathValue("user_id")
	idAsInt, err := strconv.Atoi(userId)
	if (idAsInt <= 0) || (err != nil) {
		fmt.Println("[controller] Invalid ID.", err.Error())
		utils.SendErrorResponse(respWr, 404, "Not Found")
		return
	}

	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Get Failed")
		return
	}

	fmt.Println("[controller] Current User is:", currentUser.Role)
	if currentUser.Role != models.UserRole_Admin {
		if currentUser.Id != idAsInt {
			fmt.Println("[controller] Unauthorized request on Other User")
			utils.SendErrorResponse(respWr, 401, "Unauthorized")
			return
		}
	}

	updateInput := make(map[string]interface{})
	err = utils.ExtractReqBody(req, &updateInput)
	if err != nil {
		fmt.Println("[controller]", err.Error())
		utils.SendErrorResponse(respWr, 400, "Update Failed")
		return
	}

	updateResult, updateError := services.UpdateUser(idAsInt, updateInput)
	if updateError == nil {
		utils.SendSuccessResponse(respWr, updateResult)

	} else if utils.IsErrorOfType(updateError, utils.ErrorType_UserDoesNotExist) {
		utils.SendErrorResponse(respWr, 404, "Not Found")

	} else {
		utils.SendErrorResponse(respWr, 500, "Update Failed")
	}
}

// Delete an User by ID
func DeleteUserHandler(respWr http.ResponseWriter, req *http.Request) {
	userId := req.PathValue("user_id")
	idAsInt, err := strconv.Atoi(userId)
	if (idAsInt <= 0) || (err != nil) {
		fmt.Println("[controller] Invalid ID.", err.Error())
		utils.SendErrorResponse(respWr, 404, "Not Found")
		return
	}

	currentUser, crtUserErr := middlewares.GetCurrentUser(req)
	if utils.IsErrorOfType(crtUserErr, utils.ErrorType_FormatError) {
		utils.SendErrorResponse(respWr, 400, "Get Failed")
		return
	}

	fmt.Println("[controller] Current User is:", currentUser.Role)
	if currentUser.Role != models.UserRole_Admin {
		if currentUser.Id != idAsInt {
			fmt.Println("[controller] Unauthorized request on Other User")
			utils.SendErrorResponse(respWr, 401, "Unauthorized")
			return
		}
	}

	deleteResult, deleteError := services.DeleteUser(idAsInt)
	if deleteError == nil {
		utils.SendSuccessResponse(respWr, deleteResult)

	} else if utils.IsErrorOfType(deleteError, utils.ErrorType_UserDoesNotExist) {
		utils.SendErrorResponse(respWr, 404, "Not Found")

	} else {
		utils.SendErrorResponse(respWr, 500, "Delete Failed")
	}
}
