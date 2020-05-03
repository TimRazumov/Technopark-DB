package usecase

import (
	"net/http"

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

func (useCase *UseCase) GetUsersBySlug(slug string, queryString models.QueryString) *models.Users {
	if useCase.GetBySlug(slug) == nil {
		return nil
	}
	return useCase.forumRepo.GetUsersBySlug(slug, queryString)
}

func (useCase *UseCase) GetThreadsBySlug(slug string, queryString models.QueryString) *models.Threads {
	if useCase.GetBySlug(slug) == nil {
		return nil
	}
	return useCase.forumRepo.GetThreadsBySlug(slug, queryString)
}
