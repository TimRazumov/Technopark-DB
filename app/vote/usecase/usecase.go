package usecase

import (
	"github.com/TimRazumov/Technopark-DB/app/forum"
	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/user"
)

type UseCase struct {
	userRepo  user.Repository
	forumRepo forum.Repository
}

func CreateUseCase(userRepo user.Repository, forumRepo forum.Repository) forum.UseCase {
	return &UseCase{userRepo: userRepo, forumRepo: forumRepo}
}

func (useCase *UseCase) Create(frm *models.Forum) *models.Error {
	return useCase.forumRepo.Create(*frm)
}
