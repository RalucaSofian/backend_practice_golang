package types

import (
	"app/models"
	"app/utils/date"
	"time"
)

type baseFoster struct {
	UserId      *int           `json:"userId"`
	Description *string        `json:"description"`
	PetId       *int           `json:"petId"`
	StartDate   date.DateOnly  `json:"startDate"`
	EndDate     *date.DateOnly `json:"endDate"`
}

// Foster DTO
type FosterDTO struct {
	Id int `json:"id"`
	baseFoster
	User *AuthUserDTO `json:"user"`
	Pet  *PetDTO      `json:"pet"`
}

// Create Foster Input DTO
type CreateFosterInputDTO struct {
	baseFoster
}

// New Foster (DB model) From Create Foster Input DTO
func NewFosterFromCreateFosterInput(createFosterInput CreateFosterInputDTO) models.Foster {
	var endDate *time.Time = nil
	if createFosterInput.EndDate != nil {
		endDate = &createFosterInput.EndDate.Time
	}

	foster := models.Foster{
		UserId:      createFosterInput.UserId,
		Description: createFosterInput.Description,
		PetId:       createFosterInput.PetId,
		StartDate:   createFosterInput.StartDate.Time,
		EndDate:     endDate,
	}
	return foster
}

// New Foster DTO from DB Foster
func NewFosterDtoFromDbFoster(dbFoster models.Foster) FosterDTO {
	var userDTO *AuthUserDTO = nil
	if dbFoster.User != nil {
		user := NewAuthUserDtoFromDbUser(*dbFoster.User)
		userDTO = &user
	}

	var petDTO *PetDTO = nil
	if dbFoster.Pet != nil {
		pet := NewPetDtoFromDbPet(*dbFoster.Pet)
		petDTO = &pet
	}

	startDate := date.NewDateOnlyFromTime(dbFoster.StartDate)
	var endDate *date.DateOnly = nil
	if dbFoster.EndDate != nil {
		tempEndDate := date.NewDateOnlyFromTime(*dbFoster.EndDate)
		endDate = &tempEndDate
	}

	fosterDTO := FosterDTO{
		Id: dbFoster.Id,
		baseFoster: baseFoster{
			UserId:      dbFoster.UserId,
			Description: dbFoster.Description,
			PetId:       dbFoster.PetId,
			StartDate:   startDate,
			EndDate:     endDate,
		},
		User: userDTO,
		Pet:  petDTO,
	}
	return fosterDTO
}
