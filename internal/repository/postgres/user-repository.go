package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"
	"errors"

	"gorm.io/gorm"
)

type userRepository struct {
	masterDb *gorm.DB
	clogger  *helpers.CustomLogger
}

func NewUserRepository(posgres *gorm.DB, clogger *helpers.CustomLogger) UserRepository {
	return &userRepository{masterDb: posgres, clogger: clogger}
}

type UserRepository interface {
	Tx(ctx context.Context) *gorm.DB
	SelectUserByEmailOrPhoneNumber(ctx context.Context, emailOrNo string) (*entity.User, error)
	InsertUser(ctx context.Context, tx *gorm.DB, user *entity.User) error
	SelectUserByStruct(ctx context.Context, user *entity.User) (*[]entity.User, error)
	SelectUserByStructOne(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, tx *gorm.DB, user *entity.User, newUser *entity.User) error
	SelectUserByUUID(ctx context.Context, uuid string) (*entity.User, error)
	DeleteUser(ctx context.Context, tx *gorm.DB, user *entity.User) error
	SelectUserByID(ctx context.Context, id int64) (*entity.User, error)
}

// SelectUserByStructOne implements UserRepository.
func (repo *userRepository) SelectUserByStructOne(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := repo.masterDb.Where(user).First(user).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "SelectUserByStruct.gorm.DB", err)
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) Tx(ctx context.Context) *gorm.DB {
	return repo.masterDb.Begin()
}

func (repo *userRepository) SelectUserByEmailOrPhoneNumber(ctx context.Context, emailOrNo string) (*entity.User, error) {
	var user entity.User

	err := repo.masterDb.Table(`users`).Where("email = ?", emailOrNo).Or("phone_number = ?", emailOrNo).First(&user).Error

	if err != nil {

		repo.clogger.ErrorLogger(ctx, "SelectUserByEmailOrPhoneNumber.gorm.DB", err)

		return nil, err
	}

	return &user, err
}

func (repo *userRepository) InsertUser(ctx context.Context, tx *gorm.DB, user *entity.User) error {
	err := tx.Create(&user).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "InsertUserDetail.gorm.DB", err)
	}
	return err
}

func (repo *userRepository) SelectUserByStruct(ctx context.Context, user *entity.User) (*[]entity.User, error) {
	var users []entity.User
	err := repo.masterDb.Find(&users, user).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "SelectUserByStruct.gorm.DB", err)
	}
	return &users, err
}

func (repo *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, currentData *entity.User, NewData *entity.User) error {
	if currentData == nil || NewData == nil {
		err := errors.New("nulll parameter")
		repo.clogger.ErrorLogger(ctx, "UpdateUser.gorm.DB", err)
		return err
	}
	err := tx.Model(currentData).Updates(NewData).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "UpdateUser.gorm.DB", err)
	}
	return err
}

func (repo *userRepository) SelectUserByUUID(ctx context.Context, uuid string) (*entity.User, error) {
	var user entity.User
	err := repo.masterDb.Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "SelectUserByUUID.gorm.DB", err)
		return nil, err
	}
	return &user, nil
}
func (repo *userRepository) SelectUserByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	err := repo.masterDb.Where("id = ?", id).First(&user).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "SelectUserByID.gorm.DB", err)
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) DeleteUser(ctx context.Context, tx *gorm.DB, user *entity.User) error {
	err := tx.Delete(user).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "DeleteUser.gorm.DB", err)
		return err
	}
	return nil

}
