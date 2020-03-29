package usecase

import (
	"net/http"
	"strconv"

	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/post"
	"github.com/TimRazumov/Technopark-DB/app/thread"
)

type UseCase struct {
	threadRepo thread.Repository
	postRepo   post.Repository
}

func CreateUseCase(threadRepo thread.Repository, postRepo post.Repository) post.UseCase {
	return &UseCase{threadRepo: threadRepo, postRepo: postRepo}
}

func (useCase *UseCase) Create(thrdKey string, posts *[]models.Post) *models.Error {
	if posts == nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	if len(*posts) == 0 {
		return nil
	}
	var thrd *models.Thread
	if id, err := strconv.Atoi(thrdKey); err == nil {
		thrd = useCase.threadRepo.GetByID(id)
	} else {
		thrd = useCase.threadRepo.GetBySlug(thrdKey)
	}
	if thrd == nil {
		return &models.Error{Code: http.StatusNotFound}
	}
	return useCase.postRepo.Create(*thrd, posts)
}
