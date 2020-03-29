package usecase

import (
	"github.com/TimRazumov/Technopark-DB/app/forum"
	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/user"
	"net/http"
)

type UseCase struct {
	userRepo  user.Repository
	forumRepo forum.Repository
}

func CreateUseCase(userRepo user.Repository, forumRepo forum.Repository) forum.UseCase {
	return &UseCase{userRepo: userRepo, forumRepo: forumRepo}
}

func (useCase *UseCase) Create(frm *models.Forum) *models.Error {
	if frm == nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	usr := useCase.userRepo.GetByNickName(frm.User)
	if usr == nil {
		return models.CreateNotFoundUser(frm.User)
	}
	frm.User = usr.NickName // не совпадает регистр букв
	return useCase.forumRepo.Create(*frm)
}

func (useCase *UseCase) GetBySlug(slug string) *models.Forum {
	return useCase.forumRepo.GetBySlug(slug)
}
