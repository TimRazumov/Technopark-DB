package usecase

import (
	"github.com/TimRazumov/Technopark-DB/app/forum"
	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/thread"
	"github.com/TimRazumov/Technopark-DB/app/user"
	"net/http"
	"strconv"
)

type UseCase struct {
	userRepo   user.Repository
	forumRepo  forum.Repository
	threadRepo thread.Repository
}

func CreateUseCase(userRepo user.Repository, forumRepo forum.Repository, threadRepo thread.Repository) thread.UseCase {
	return &UseCase{userRepo: userRepo, forumRepo: forumRepo, threadRepo: threadRepo}
}

func (useCase *UseCase) Create(thrd *models.Thread) *models.Error {
	if thrd == nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	frm := useCase.forumRepo.GetBySlug(thrd.Forum)
	if frm == nil {
		return models.CreateNotFoundForumThread(thrd.Forum)
	}
	thrd.Forum = frm.Slug
	usr := useCase.userRepo.GetByNickName(thrd.Author)
	if usr == nil {
		return models.CreateNotFoundAuthorThread(thrd.Author)
	}
	thrd.Author = usr.NickName
	return useCase.threadRepo.Create(thrd)
}

func (useCase *UseCase) GetByID(id int) *models.Thread {
	return useCase.threadRepo.GetByID(id)
}

func (useCase *UseCase) GetBySlug(slug string) *models.Thread {
	return useCase.threadRepo.GetBySlug(slug)
}

func (useCase *UseCase) Update(newThrd *models.Thread) *models.Error {
	return useCase.threadRepo.Update(newThrd)
}

func (useCase *UseCase) UpdateVote(thrdKey string, vt models.Vote) *models.Thread {
	var thrd *models.Thread
	if id, err := strconv.Atoi(thrdKey); err == nil {
		thrd = useCase.threadRepo.GetByID(id)
	} else {
		thrd = useCase.threadRepo.GetBySlug(thrdKey)
	}
	if thrd == nil {
		return nil
	}
	vt.Thread = thrd.ID
	err := useCase.threadRepo.UpdateVote(vt)
	if err != nil {
		return nil
	}
	return useCase.threadRepo.GetByID(thrd.ID)
}

func (useCase *UseCase) GetPostsBySlugOrID(thrdKey string, queryString models.QueryString) *[]models.Post {
	var thrd *models.Thread
	if id, err := strconv.Atoi(thrdKey); err == nil {
		thrd = useCase.GetByID(id)
	} else {
		thrd = useCase.GetBySlug(thrdKey)
	}
	if thrd == nil {
		return nil
	}
	return useCase.threadRepo.GetPostsByThread(*thrd, queryString)
}
