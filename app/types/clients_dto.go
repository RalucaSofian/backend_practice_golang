package types

import (
	"app/models"
)

type baseClient struct {
	UserId      *int    `json:"userId"`
	Description *string `json:"description"`
}

// Client DTO
type ClientDTO struct {
	Id int `json:"id"`
	baseClient
	User *AuthUserDTO `json:"user"`
}

// Create Client Input DTO
type CreateClientInputDTO struct {
	baseClient
}

// New Client (DB model) from Create Client Input DTO
func NewClientFromCreateClientInput(createClientInput CreateClientInputDTO) models.Client {
	client := models.Client{
		UserId:      createClientInput.UserId,
		Description: createClientInput.Description,
	}
	return client
}

// New Client DTO from DB Client
func NewClientDtoFromDbClient(dbClient models.Client) ClientDTO {
	var userDTO *AuthUserDTO = nil
	if dbClient.User != nil {
		user := NewAuthUserDtoFromDbUser(*dbClient.User)
		userDTO = &user
	}

	clientDTO := ClientDTO{
		baseClient: baseClient{
			UserId:      dbClient.UserId,
			Description: dbClient.Description,
		},
		Id:   dbClient.Id,
		User: userDTO,
	}
	return clientDTO
}
