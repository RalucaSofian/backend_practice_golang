package types

import "app/models"

// Possible Pet Genders
type PetGender string

const (
	PetGender_Male   PetGender = "MALE"
	PetGender_Female PetGender = "FEMALE"
)

type basePet struct {
	Name        string     `json:"name"`
	Species     *string    `json:"species"`
	Gender      *PetGender `json:"gender"`
	Age         *float32   `json:"age"`
	Description *string    `json:"description"`
}

// Pet DTO
type PetDTO struct {
	basePet
	Id int `json:"id"`
}

// Create Pet Input DTO
type CreatePetInputDTO struct {
	basePet
}

// New Pet (DB model) from Create Pet Input DTO
func NewPetFromCreatePetInput(createPetInput CreatePetInputDTO) models.Pet {
	pet := models.Pet{
		Name:        createPetInput.Name,
		Species:     createPetInput.Species,
		Gender:      (*string)(createPetInput.Gender),
		Age:         createPetInput.Age,
		Description: createPetInput.Description,
	}
	return pet
}

// New Pet DTO from Db Pet
func NewPetDtoFromDbPet(dbPet models.Pet) PetDTO {
	petDTO := PetDTO{
		basePet: basePet{
			Name:        dbPet.Name,
			Species:     dbPet.Species,
			Gender:      (*PetGender)(dbPet.Gender),
			Age:         dbPet.Age,
			Description: dbPet.Description,
		},
		Id: dbPet.Id,
	}
	return petDTO
}
