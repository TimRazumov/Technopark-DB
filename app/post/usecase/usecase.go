package usecase

import (
	"net/http"
	"strconv"

	"github.com/TimRazumov/Technopark-DB/app/forum"
	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/post"
	"github.com/TimRazumov/Technopark-DB/app/thread"
	"github.com/TimRazumov/Technopark-DB/app/user"
)

type UseCase struct {
	userRepo   user.Repository
	forumRepo  forum.Repository
	threadRepo thread.Repository
	postRepo   post.Repository
}

func CreateUseCase(userRepo user.Repository, forumRepo forum.Repository,
	threadRepo thread.Repository, postRepo post.Repository) post.UseCase {
	return &UseCase{
		userRepo:   userRepo,
		forumRepo:  forumRepo,
		threadRepo: threadRepo,
		postRepo:   postRepo}
}

func (useCase *UseCase) Create(thrdKey string, posts *models.Posts) *models.Error {
	if posts == nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	var thrd *models.Thread
	id, err := strconv.Atoi(thrdKey)
	if err == nil {
		thrd = useCase.threadRepo.GetByID(id)
	} else {
		thrd = useCase.threadRepo.GetBySlug(thrdKey)
	}
	if thrd == nil {
		return models.CreateNotFoundThreadPost(id)
	}
	return useCase.postRepo.Create(*thrd, posts)
}

func (useCase *UseCase) GetByID(id int, options models.Related) *models.PostGet {
	var postAll models.PostGet
	postAll.Post = useCase.postRepo.GetByID(id)
	if postAll.Post == nil {
		return nil
	}
	pst := *postAll.Post
	if options.User {
		postAll.Author = useCase.userRepo.GetByNickName(pst.Author)
		if postAll.Author == nil {
			return nil
		}
	}
	if options.Forum {
		postAll.Forum = useCase.forumRepo.GetBySlug(pst.Forum)
		if postAll.Forum == nil {
			return nil
		}
	}
	if options.Thread {
		postAll.Thread = useCase.threadRepo.GetByID(pst.Thread)
		if postAll.Thread == nil {
			return nil
		}
	}
	return &postAll
}

func (useCase *UseCase) Update(pst *models.Post) *models.Error {
	return useCase.postRepo.Update(pst)
}
