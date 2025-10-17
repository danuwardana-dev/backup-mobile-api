package postgres

import (
	"backend-mobile-api/helpers"
	"backend-mobile-api/model/entity"
	"context"
	"errors"
	"gorm.io/gorm"
)

type userDetailRepository struct {
	masterDb *gorm.DB
	clogger  *helpers.CustomLogger
}

func NewUserDetailRepository(posgres *gorm.DB, clogger *helpers.CustomLogger) UserDetailRepository {
	return &userDetailRepository{masterDb: posgres, clogger: clogger}
}

type UserDetailRepository interface {
	Tx(ctx context.Context) *gorm.DB
	InsertUserDetail(ctx context.Context, tx *gorm.DB, user *entity.UserDetail) error
	SelectUserDetailByEmailOrPhoneNumber(ctx context.Context, emailOrNo string) (*entity.UserDetail, error)
	UpdateUserDetail(ctx context.Context, tx *gorm.DB, curentUserDetail *entity.UserDetail, newUserDetail *entity.UserDetail) error
	SelectUserDetailByUserId(ctx context.Context, userId int64) (*entity.UserDetail, error)
	SelectUserDetailByUserUUID(ctx context.Context, userUUID *string) (*entity.UserDetail, error)
}

func (ur *userDetailRepository) Tx(ctx context.Context) *gorm.DB {
	return ur.masterDb.WithContext(ctx)
}
func (repo *userDetailRepository) InsertUserDetail(ctx context.Context, tx *gorm.DB, user *entity.UserDetail) error {
	err := tx.Create(user).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "InsertUserDetail.tx.Create(user).Error", err)
	}
	return err
}
func (repo *userDetailRepository) SelectUserDetailByEmailOrPhoneNumber(ctx context.Context, emailOrNo string) (*entity.UserDetail, error) {
	var userDetail entity.UserDetail
	err := repo.masterDb.Table(`user_details`).Where("email = ?", emailOrNo).Or("phone_number = ?", emailOrNo).First(&userDetail).Error

	if err != nil {

		repo.clogger.ErrorLogger(ctx, "SelectUserDetailByEmailOrPhoneNumber.gorm.DB", err)

		return nil, err
	}
	return &userDetail, nil
}
func (repo *userDetailRepository) UpdateUserDetail(ctx context.Context, tx *gorm.DB, curentUserDetail *entity.UserDetail, newUserDetail *entity.UserDetail) error {

	if curentUserDetail == nil || newUserDetail == nil {
		err := errors.New("nulll parameter")
		repo.clogger.ErrorLogger(ctx, "UpdateUserDetail.gorm.DB", err)
		return err
	}
	err := tx.Model(curentUserDetail).Updates(newUserDetail).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "UpdateUser.gorm.DB", err)
	}
	return err

}
func (repo *userDetailRepository) SelectUserDetailByUserId(ctx context.Context, userId int64) (*entity.UserDetail, error) {
	var userDetail entity.UserDetail
	err := repo.masterDb.WithContext(ctx).Where("user_id = ?", userId).First(&userDetail).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "SelectUserDetailByUserId.gorm.DB", err)
	}
	return &userDetail, err
}
func (repo *userDetailRepository) SelectUserDetailByUserUUID(ctx context.Context, userUUID *string) (*entity.UserDetail, error) {
	var userDetail entity.UserDetail
	err := repo.masterDb.WithContext(ctx).Where("user_uuid = ?", *userUUID).First(&userDetail).Error
	if err != nil {
		repo.clogger.ErrorLogger(ctx, "SelectUserDetailByUserUUID.gorm.DB", err)
	}
	return &userDetail, err
}
