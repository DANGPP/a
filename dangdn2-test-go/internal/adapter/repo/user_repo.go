package repo

import (
	"dangdn2-test-go/internal/core/domain"

	"gorm.io/gorm"
)

type GormUserRepo struct {
	DB *gorm.DB
}

func (r *GormUserRepo) GetAll() ([]domain.User, error) {
	var users []domain.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *GormUserRepo) Create(user domain.User) error {
	return r.DB.Create(&user).Error
}

func (r *GormUserRepo) FindByNameAndPass(name, pass string) (domain.User, error) {
	var user domain.User
	err := r.DB.Where("name = ? AND pass = ?", name, pass).First(&user).Error
	return user, err
}
