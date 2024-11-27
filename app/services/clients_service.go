package services

import (
	"app/db"
	"app/models"
	"app/types"
	"app/utils"
	"app/utils/query_utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Create a new Client
func CreateClient(createClientInput types.CreateClientInputDTO) (*types.ClientDTO, error) {
	dbClient := types.NewClientFromCreateClientInput(createClientInput)

	_, err := db.GetConn().NewInsert().Model(&dbClient).Exec(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	err = db.GetConn().NewSelect().
		Model(&dbClient).
		Where("client.id = ?", dbClient.Id).
		Relation("User").
		Scan(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	clientDTO := types.NewClientDtoFromDbClient(dbClient)
	return &clientDTO, nil
}

// Get a Client by ID
func GetClientById(clientId int) (*types.ClientDTO, error) {
	dbClient := models.Client{}

	err := db.GetConn().NewSelect().
		Model(&dbClient).
		Where("client.id = ?", clientId).
		Relation("User").
		Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("[service] Not Found")
		return nil, utils.NewApiError(utils.ErrorType_ClientDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	clientDTO := types.NewClientDtoFromDbClient(dbClient)
	return &clientDTO, nil
}

// Get all Clients (+ Query)
func GetAllClients(queryInfo *query_utils.QueryInfo) ([]*types.ClientDTO, error) {
	dbClients := []models.Client{}

	query := db.GetConn().NewSelect().Model(&dbClients).Relation("User")
	if queryInfo != nil {
		query = queryInfo.Process(query)
	}
	queryErr := query.Scan(context.Background())
	if queryErr != nil {
		fmt.Println("[service]", queryErr.Error())
		return nil, queryErr
	}

	clients := []*types.ClientDTO{}
	for _, dbClient := range dbClients {
		clientDTO := types.NewClientDtoFromDbClient(dbClient)
		clients = append(clients, &clientDTO)
	}

	return clients, nil
}

// Update a Client by ID
func UpdateClient(clientId int, updateInput map[string]interface{}) (*types.ClientDTO, error) {
	dbClient := models.Client{}

	err := db.GetConn().NewSelect().
		Model(&dbClient).
		Where("client.id = ?", clientId).
		Relation("User").
		Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("[service] Not Found")
		return nil, utils.NewApiError(utils.ErrorType_ClientDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	dbClient, err = models.Update(dbClient, updateInput)
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	_, err = db.GetConn().NewUpdate().
		Model(&dbClient).
		WherePK().
		Exec(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	clientDTO := types.NewClientDtoFromDbClient(dbClient)
	return &clientDTO, nil
}

// Delete a Client by ID
func DeleteClient(clientId int) (string, error) {
	dbClient := models.Client{}

	exists, err := db.GetConn().NewSelect().
		Model(&dbClient).
		Where("id = ?", clientId).
		Exists(context.Background())
	if !exists {
		return "Delete Failed", utils.NewApiError(utils.ErrorType_ClientDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err.Error())
		return "Delete Failed", err
	}

	_, err = db.GetConn().NewDelete().
		Model(&dbClient).
		Where("id = ?", clientId).
		Exec(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return "Delete Failed", err
	}

	return "Delete Success", nil
}
