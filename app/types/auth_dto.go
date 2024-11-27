package types

import "app/models"

type baseUser struct {
	Email   string  `json:"email"`
	Name    *string `json:"name"`
	Address *string `json:"address"`
	Phone   *string `json:"phone"`
	Role    string  `json:"role"`
}

// Auth User DTO
type AuthUserDTO struct {
	Id int `json:"id"`
	baseUser
}

// Register Input DTO
type RegisterInputDTO struct {
	baseUser
	Password string `json:"password"`
}

// Login Input DTO
type LoginInputDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Update User Input DTO
type UpdateUserInputDTO struct {
	Email   *string `json:"email"`
	Name    *string `json:"name"`
	Address *string `json:"address"`
	Phone   *string `json:"phone"`
}

// New Auth User (DB model) From Register Input
func NewAuthUserFromRegisterInput(registerInput RegisterInputDTO) models.AuthUser {
	authUser := models.AuthUser{
		Email:    registerInput.Email,
		Password: registerInput.Password,
		Name:     registerInput.Name,
		Address:  registerInput.Address,
		Phone:    registerInput.Phone,
	}
	return authUser
}

// New Auth User DTO From DB User
func NewAuthUserDtoFromDbUser(dbUser models.AuthUser) AuthUserDTO {
	authUserDTO := AuthUserDTO{
		baseUser: baseUser{
			Email:   dbUser.Email,
			Name:    dbUser.Name,
			Address: dbUser.Address,
			Phone:   dbUser.Phone,
			Role:    string(dbUser.Role),
		},
		Id: dbUser.Id,
	}
	return authUserDTO
}

// New Auth User (DB model) From Update Input
func NewAuthuserFromUpdateInput(updateInput UpdateUserInputDTO) models.AuthUser {
	authUser := models.AuthUser{
		Email:   *updateInput.Email,
		Name:    updateInput.Name,
		Address: updateInput.Address,
		Phone:   updateInput.Phone,
	}
	return authUser
}
