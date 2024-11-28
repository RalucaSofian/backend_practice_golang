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

// Create a new Foster
func CreateFoster(createFosterInput types.CreateFosterInputDTO) (*types.FosterDTO, error) {
	dbFoster := types.NewFosterFromCreateFosterInput(createFosterInput)

	_, err := db.GetConn().NewInsert().Model(&dbFoster).Exec(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	err = db.GetConn().NewSelect().
		Model(&dbFoster).
		Where("foster.id = ?", dbFoster.Id).
		Relation("User").
		Relation("Pet").
		Scan(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	fosterDTO := types.NewFosterDtoFromDbFoster(dbFoster)
	return &fosterDTO, nil
}

// Get a Foster by ID
func GetFosterById(fosterId int) (*types.FosterDTO, error) {
	dbFoster := models.Foster{}

	err := db.GetConn().NewSelect().
		Model(&dbFoster).
		Where("foster.id = ?", fosterId).
		Relation("User").
		Relation("Pet").
		Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("[service] Not Found")
		return nil, utils.NewApiError(utils.ErrorType_FosterDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	fosterDTO := types.NewFosterDtoFromDbFoster(dbFoster)
	return &fosterDTO, nil
}

// Get all Fosters (+ Query)
func GetAllFosters(queryInfo *query_utils.QueryInfo) ([]*types.FosterDTO, error) {
	dbFosters := []models.Foster{}

	query := db.GetConn().NewSelect().Model(&dbFosters)
	if queryInfo != nil {
		query = queryInfo.Process(query)
	}
	queryErr := query.Relation("User").Relation("Pet").Scan(context.Background())
	if queryErr != nil {
		fmt.Println("[service]", queryErr.Error())
		return nil, queryErr
	}

	fosters := []*types.FosterDTO{}
	for _, dbFoster := range dbFosters {
		fosterDTO := types.NewFosterDtoFromDbFoster(dbFoster)
		fosters = append(fosters, &fosterDTO)
	}

	return fosters, nil
}

// Update a Foster by ID
func UpdateFoster(fosterId int, updateInput map[string]interface{}) (*types.FosterDTO, error) {
	dbFoster := models.Foster{}

	err := db.GetConn().NewSelect().
		Model(&dbFoster).
		Where("foster.id = ?", fosterId).
		Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("[service] Not Found")
		return nil, utils.NewApiError(utils.ErrorType_FosterDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	dbFoster, err = models.Update(dbFoster, updateInput)
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	_, err = db.GetConn().NewUpdate().
		Model(&dbFoster).
		WherePK().
		Exec(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	// New Select for the complete updated data
	err = db.GetConn().NewSelect().
		Model(&dbFoster).
		Where("foster.id = ?", fosterId).
		Relation("User").
		Relation("Pet").
		Scan(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return nil, err
	}

	fosterDTO := types.NewFosterDtoFromDbFoster(dbFoster)
	return &fosterDTO, nil
}

// Delete a Foster by ID
func DeleteFoster(fosterId int) (string, error) {
	dbFoster := models.Foster{}

	exists, err := db.GetConn().NewSelect().
		Model(&dbFoster).
		Where("foster.id = ?", fosterId).
		Exists(context.Background())
	if !exists {
		return "Delete Failed", utils.NewApiError(utils.ErrorType_FosterDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err.Error())
		return "Delete Failed", err
	}

	_, err = db.GetConn().NewDelete().
		Model(&dbFoster).
		Where("foster.id = ?", fosterId).
		Exec(context.Background())
	if err != nil {
		fmt.Println("[service]", err.Error())
		return "Delete Failed", err
	}

	return "Delete Success", nil
}
