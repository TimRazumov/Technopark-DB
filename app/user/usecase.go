package user

import "github.com/TimRazumov/Technopark-DB/app/models"

type UseCase interface {
	Create(usr models.User) *models.Error
	GetByNickName(nickName string) *models.User
	GetByEmail(email string) *models.User
	Update(newUsr *models.User) *models.Error
}
