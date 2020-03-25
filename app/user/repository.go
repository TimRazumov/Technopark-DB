package user

import "github.com/TimRazumov/Technopark-DB/app/models"

type Repository interface {
	Create(usr models.User) *models.Error
	GetByNickName(nickname string) *models.User
	GetByEmail(email string) *models.User
	Update(newUsr *models.User) *models.Error
}
