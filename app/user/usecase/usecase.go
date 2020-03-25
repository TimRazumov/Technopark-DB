package usecase

import (
	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/user"
)

type UseCase struct {
	userRepo user.Repository
}

func CreateUseCase(userRepo user.Repository) user.UseCase {
	return &UseCase{userRepo: userRepo}
}

func (useCase *UseCase) Create(usr models.User) *models.Error {
	return useCase.userRepo.Create(usr)
}

func (useCase *UseCase) GetByNickName(nickName string) *models.User {
	return useCase.userRepo.GetByNickName(nickName)
}

func (useCase *UseCase) GetByEmail(email string) *models.User {
	return useCase.userRepo.GetByEmail(email)
}

func (useCase *UseCase) Update(newUsr *models.User) *models.Error {
	return useCase.userRepo.Update(newUsr)
}
