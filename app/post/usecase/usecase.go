package usecase

import (
	"net/http"

	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/post"
)

type UseCase struct {
	postRepo post.Repository
}

func CreateUseCase(postRepo post.Repository) post.UseCase {
	return &UseCase{postRepo: postRepo}
}

func (useCase *UseCase) Create(postKey string, posts *[]models.Post) *models.Error {
	if posts == nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}

	/*_, err := fmt.Sscan(postKey, &post.ID)
	if err != nil {

	}*/
}
