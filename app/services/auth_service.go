package services

import (
	"app/db"
	"app/models"
	"app/types"
	"app/utils"
	"app/utils/access_utils"
	"app/utils/query_utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Register a new Auth User
func RegisterUser(registerInput types.RegisterInputDTO) (*types.AuthUserDTO, error) {
	dbUser := types.NewAuthUserFromRegisterInput(registerInput)

	exists, err := db.GetConn().NewSelect().
		Model(&dbUser).
		Where("email = ?", dbUser.Email).
		Exists(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}
	if exists {
		fmt.Println("[service] User Already Exists")
		return nil, utils.NewApiError(utils.ErrorType_UserAlreadyExists, "User Already Exists")
	}

	hashedPassword, err := access_utils.GeneratePasswordHash(registerInput.Password)
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}
	dbUser.Password = hashedPassword
	dbUser.Role = models.UserRole_User

	_, err = db.GetConn().NewInsert().Model(&dbUser).Exec(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}
	userDTO := types.NewAuthUserDtoFromDbUser(dbUser)

	// Also create a Client object for the Current User
	newClientInput := types.CreateClientInputDTO{}
	clientDescription := fmt.Sprintf("%s (%s)", *dbUser.Name, *dbUser.Phone)
	newClientInput.UserId = &dbUser.Id
	newClientInput.Description = &clientDescription

	newClientResult, newClientErr := CreateClient(newClientInput)
	fmt.Println("[service] Created Client object for the current Auth User. Result:", newClientResult, "Error:", newClientErr)

	return &userDTO, nil
}

// Login as an existing Auth User
func LoginUser(loginInput types.LoginInputDTO) (string, error) {
	userModel := models.AuthUser{}
	loginToken := ""

	err := db.GetConn().NewSelect().
		Model(&userModel).
		Where("email = ?", loginInput.Email).
		Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("[service] Not Found")
		return "", utils.NewApiError(utils.ErrorType_UserLoginFailed, "Login Failed")
	}
	if err != nil {
		fmt.Println("[service]", err.Error())
		return "", err
	}

	validPassword := access_utils.VerifyPasswordHash(loginInput.Password, userModel.Password)
	if !validPassword {
		return "", utils.NewApiError(utils.ErrorType_UserLoginFailed, "Login Failed")
	}

	//! can return ErrorType_JWTError
	loginToken, err = access_utils.CreateAccessToken(loginInput.Email)
	if err != nil {
		return "", utils.NewApiError(utils.ErrorType_UserLoginFailed, "Login Failed")
	}

	return loginToken, err
}

// Get an Auth User by ID
func GetUserById(userId int) (*types.AuthUserDTO, error) {
	dbUser := models.AuthUser{}

	err := db.GetConn().NewSelect().
		Model(&dbUser).
		Where("id = ?", userId).
		Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("[service] Not Found")
		return nil, utils.NewApiError(utils.ErrorType_UserDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	authUser := types.NewAuthUserDtoFromDbUser(dbUser)
	return &authUser, nil
}

// Get an Auth User by Email
func GetUserByEmail(email string) (*models.AuthUser, error) {
	dbUser := models.AuthUser{}

	err := db.GetConn().NewSelect().
		Model(&dbUser).
		Where("email = ?", email).
		Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("[service] Not Found")
		return nil, utils.NewApiError(utils.ErrorType_UserDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	return &dbUser, nil
}

// Get all Auth Users (+ Query)
func GetAllUsers(queryInfo *query_utils.QueryInfo) ([]*types.AuthUserDTO, error) {
	dbUsers := []models.AuthUser{}

	query := db.GetConn().NewSelect().Model(&dbUsers)
	if queryInfo != nil {
		query = queryInfo.Process(query)
	}
	queryErr := query.Scan(context.Background())
	if queryErr != nil {
		fmt.Println("[service]", queryErr.Error())
		return nil, queryErr
	}

	authUsers := []*types.AuthUserDTO{}
	for _, dbUser := range dbUsers {
		userDto := types.NewAuthUserDtoFromDbUser(dbUser)
		authUsers = append(authUsers, &userDto)
	}

	return authUsers, nil
}

// Update an Auth User by ID
func UpdateUser(userId int, updateInput map[string]interface{}) (*types.AuthUserDTO, error) {
	dbUser := models.AuthUser{}

	err := db.GetConn().NewSelect().
		Model(&dbUser).
		Where("id = ?", userId).
		Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("[service] Not Found")
		return nil, utils.NewApiError(utils.ErrorType_UserDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	dbUser, err = models.Update(dbUser, updateInput)
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	_, err = db.GetConn().NewUpdate().
		Model(&dbUser).
		WherePK().
		Exec(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	authUser := types.NewAuthUserDtoFromDbUser(dbUser)
	return &authUser, nil
}

// Delete an Auth User by ID
func DeleteUser(userId int) (string, error) {
	dbUser := models.AuthUser{}

	exists, err := db.GetConn().NewSelect().
		Model(&dbUser).
		Where("id = ?", userId).
		Exists(context.Background())
	if !exists {
		return "Delete Failed", utils.NewApiError(utils.ErrorType_UserDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err.Error())
		return "Delete Failed", err
	}

	_, err = db.GetConn().NewDelete().
		Model(&dbUser).
		Where("id = ?", userId).
		Exec(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return "Delete Failed", err
	}

	return "Delete Success", nil
}
