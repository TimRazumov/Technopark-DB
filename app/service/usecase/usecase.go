package usecase

import (
	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/service"
)

type UseCase struct {
	serviceRepo service.Repository
}

func CreateUseCase(serviceRepo service.Repository) service.UseCase {
	return &UseCase{serviceRepo: serviceRepo}
}

func (useCase *UseCase) Get() *models.Status {
	return useCase.serviceRepo.Get()
}

func (useCase *UseCase) Clear() *models.Error {
	return useCase.serviceRepo.Clear()
}
