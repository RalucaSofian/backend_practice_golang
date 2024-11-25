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

func CreatePet(createPetInput types.CreatePetInputDTO) (*types.PetDTO, error) {
	dbPet := types.NewPetFromCreatePetInput(createPetInput)

	_, err := db.GetConn().NewInsert().Model(&dbPet).Exec(context.Background())
	if err != nil {
		fmt.Println("[service]", err)
		return nil, err
	}

	petDTO := types.NewPetDtoFromDbPet(dbPet)
	return &petDTO, nil
}

// Get a Pet by ID
func GetPetById(petId int) (*types.PetDTO, error) {
	dbPet := models.Pet{}

	err := db.GetConn().NewSelect().
		Model(&dbPet).
		Where("id = ?", petId).
		Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("[service] Not Found")
		return nil, utils.NewApiError(utils.ErrorType_PetDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err)
		return nil, err
	}

	petDTO := types.NewPetDtoFromDbPet(dbPet)
	return &petDTO, nil
}

// Get all Pets (+ Query)
func GetAllPets(queryInfo *query_utils.QueryInfo) ([]*types.PetDTO, error) {
	dbPets := []models.Pet{}

	query := db.GetConn().NewSelect().Model(&dbPets)
	if queryInfo != nil {
		query = queryInfo.Process(query)
	}
	queryErr := query.Scan(context.Background())
	if queryErr != nil {
		fmt.Println("[service]", queryErr.Error())
		return nil, queryErr
	}

	pets := []*types.PetDTO{}
	for _, dbPet := range dbPets {
		petDTO := types.NewPetDtoFromDbPet(dbPet)
		pets = append(pets, &petDTO)
	}

	return pets, nil
}

// Update a Pet by ID
func UpdatePet(petId int, updateInput map[string]interface{}) (*types.PetDTO, error) {
	dbPet := models.Pet{}

	err := db.GetConn().NewSelect().
		Model(&dbPet).
		Where("id = ?", petId).
		Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("[service] Not Found")
		return nil, utils.NewApiError(utils.ErrorType_PetDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err)
		return nil, err
	}

	dbPet, err = models.Update(dbPet, updateInput)
	if err != nil {
		fmt.Println("[service]", err)
		return nil, err
	}

	_, err = db.GetConn().NewUpdate().
		Model(&dbPet).
		WherePK().
		Exec(context.Background())
	if err != nil {
		fmt.Println("[service]", err)
		return nil, err
	}

	petDTO := types.NewPetDtoFromDbPet(dbPet)
	return &petDTO, nil
}

// Delete a Pet by ID
func DeletePet(petID int) (string, error) {
	dbPet := models.Pet{}

	exists, err := db.GetConn().NewSelect().
		Model(&dbPet).
		Where("id = ?", petID).
		Exists(context.Background())
	if !exists {
		return "Delete Failed", utils.NewApiError(utils.ErrorType_PetDoesNotExist, "Not Found")
	}
	if err != nil {
		fmt.Println("[service]", err)
		return "Delete Failed", err
	}

	_, err = db.GetConn().NewDelete().
		Model(&dbPet).
		Where("id = ?", petID).
		Exec(context.Background())
	if err != nil {
		fmt.Println("[service]", err)
		return "Delete Failed", err
	}

	return "Delete Success", nil
}
